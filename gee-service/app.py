# app.py (Final, Complete, and Corrected Version)
import ee
from flask import Flask, request, jsonify
import os
from datetime import datetime, timedelta
import requests # Make sure you have run 'pip install requests'
from flask import Response

import json # We will need this to handle the computation graph


# --- Initialization ---
app = Flask(__name__)

# --- CONFIGURATION ---
# Get your Google Cloud Project ID from the Cloud Console.
PROJECT_ID = 'geowatch-project' # This should be your correct Project ID
KEY_FILE_PATH = 'gee-credentials.json'

# --- AUTHENTICATION LOGIC ---
try:
    if not os.path.exists(KEY_FILE_PATH):
        raise FileNotFoundError(f"Service account key file not found at: {KEY_FILE_PATH}")

    credentials = ee.ServiceAccountCredentials(email=None, key_file=KEY_FILE_PATH)
    ee.Initialize(credentials=credentials, project=PROJECT_ID)
    
    print("✅ Successfully authenticated with Google Earth Engine.")

except Exception as e:
    print(f"❌ Failed to authenticate with GEE. Please check your PROJECT_ID and key file path.")
    print(f"Error details: {e}")

# In gee-service/app.py
# In gee-service/app.py

@app.route('/analyze-changes', methods=['POST'])
def analyze_changes():
    data = request.get_json(); aoi_coords = data.get('aoi'); start_date_str = data.get('startDate'); end_date_str = data.get('endDate');
    start_year = datetime.strptime(start_date_str, '%Y-%m-%d').year; end_year = datetime.strptime(end_date_str, '%Y-%m-%d').year
    
    print(f"Received PRODUCTION request for years: {start_year} vs {end_year}")

    try:
        aoi = ee.Geometry.Polygon(aoi_coords)

        # 1. --- Analysis Logic (This is perfect, no changes needed) ---
        def create_analysis_composite(year):
            current_year = datetime.now().year
            if year >= current_year: raise ValueError(f"Data for year {year} is not yet available.")
            startDate = ee.Date.fromYMD(year, 7, 1); 
            endDate = ee.Date.fromYMD(year, 7, 15);

            collection_name = ''; nir_band, red_band, green_band, swir1_band = '', '', '', ''
            if year >= 2022: collection_name, nir_band, red_band, green_band, swir1_band = 'LANDSAT/LC09/C02/T1_L2', 'SR_B5', 'SR_B4', 'SR_B3', 'SR_B6'
            elif year >= 2013: collection_name, nir_band, red_band, green_band, swir1_band = 'LANDSAT/LC08/C02/T1_L2', 'SR_B5', 'SR_B4', 'SR_B3', 'SR_B6'
            elif year >= 1999: collection_name, nir_band, red_band, green_band, swir1_band = 'LANDSAT/LE07/C02/T1_L2', 'SR_B4', 'SR_B3', 'SR_B2', 'SR_B5'
            elif year >= 1984: collection_name, nir_band, red_band, green_band, swir1_band = 'LANDSAT/LT05/C02/T1_L2', 'SR_B4', 'SR_B3', 'SR_B2', 'SR_B5'
            else: raise ValueError(f"No reliable data for year {year}.")
            def maskL8sr(image):
                qa = image.select('QA_PIXEL'); mask = qa.bitwiseAnd(1 << 3).eq(0).And(qa.bitwiseAnd(1 << 4).eq(0)); return image.updateMask(mask).select("SR_B.*").multiply(0.0000275).add(-0.2)
            collection = ee.ImageCollection(collection_name).filterBounds(aoi).filterDate(startDate, endDate)
            composite = collection.map(maskL8sr).median()
            band_names = composite.bandNames().getInfo()
            if not band_names: raise ValueError(f"No cloud-free images found for {year}.")
            ndvi = composite.normalizedDifference([nir_band, red_band]).rename('ndvi'); ndwi = composite.normalizedDifference([green_band, nir_band]).rename('ndwi'); ndbi = composite.normalizedDifference([swir1_band, nir_band]).rename('ndbi')
            return composite.addBands([ndvi, ndwi, ndbi])

        image_before = create_analysis_composite(start_year).select(['ndvi', 'ndwi', 'ndbi'])
        image_after = create_analysis_composite(end_year).select(['ndvi', 'ndwi', 'ndbi'])
        
        diff = image_after.subtract(image_before); squared_diff = diff.pow(2)
        change_magnitude = squared_diff.reduce(ee.Reducer.sum()).rename('change_sum')

        # 2. --- Statistics Logic (This is also perfect) ---
        print("Calculating statistics...")
        stats = change_magnitude.reduceRegion(
            reducer=ee.Reducer.minMax().combine(ee.Reducer.percentile([98]), sharedInputs=True),
            geometry=aoi,
            scale=90, # A reasonable scale for stats calculation
            maxPixels=1e9
        ).getInfo()
        min_val = stats['change_sum_min']
        p98_val = stats['change_sum_p98']
        print(f"Dynamic range found: min={min_val}, p98_max={p98_val}")

        # 3. --- THE FINAL, CORRECTED `computePixels` REQUEST ---
        print("Building request with Affine Transform...")
        grid_dimensions = 2048 # High resolution

        # --- Calculate the Affine Transform ---
        aoi_bounds = aoi.bounds().getInfo()['coordinates'][0]
        min_lon, min_lat = aoi_bounds[0][0], aoi_bounds[0][1]
        max_lon, max_lat = aoi_bounds[2][0], aoi_bounds[2][1]
        scale_x = (max_lon - min_lon) / grid_dimensions
        scale_y = (min_lat - max_lat) / grid_dimensions
        translate_x = min_lon
        translate_y = max_lat
        
        request_params = {
            'expression': change_magnitude.selfMask(),
            'fileFormat': 'PNG',
            'grid': {
                'dimensions': { 'width': grid_dimensions, 'height': grid_dimensions },
                'affineTransform': {
                    'scaleX': scale_x, 'shearX': 0, 'translateX': translate_x,
                    'shearY': 0, 'scaleY': scale_y, 'translateY': translate_y
                },
                'crsCode': 'EPSG:4326'
            },
            'bandIds': ['change_sum'],
            'visualizationOptions': {
                'ranges': [{'min': min_val, 'max': p98_val}], # Use a slightly smaller min to ensure visibility
                'paletteColors': ['0000FF', 'FFFF00', 'FFA500', 'FF0000']
            }
        }
        
        pixel_data = ee.data.computePixels(request_params)
        print("...PNG data received.")
        
        return Response(pixel_data, mimetype='image/png')

    except Exception as e:
        print(f"ERROR: A GEE error occurred: {e}")
        return jsonify({"error": f"An error occurred during GEE analysis: {str(e)}"}), 500

# --- Run the Server ---
if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True, port=5000)


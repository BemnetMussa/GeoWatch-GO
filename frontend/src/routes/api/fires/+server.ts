import { json } from "@sveltejs/kit"
import type { RequestHandler } from "@sveltejs/kit"
import type { FireData } from "$lib/types"

const NASA_FIRMS_BASE_URL = "https://firms.modaps.eosdis.nasa.gov/api/area/csv"
const MAP_KEY = "YOUR_NASA_FIRMS_MAP_KEY" // This should be in environment variables

export const GET: RequestHandler = async ({ url }) => {
  try {
    const bounds = url.searchParams.get("bounds")
    const days = url.searchParams.get("days") || "1"
    const source = url.searchParams.get("source") || "VIIRS_SNPP_NRT"

    if (!bounds) {
      return json({ error: "Bounds parameter is required" }, { status: 400 })
    }

    // Parse bounds: "west,south,east,north"
    const [west, south, east, north] = bounds.split(",").map(Number)

    if ([west, south, east, north].some(isNaN)) {
      return json({ error: "Invalid bounds format" }, { status: 400 })
    }

    // Construct NASA FIRMS API URL
    const firmsUrl = `${NASA_FIRMS_BASE_URL}/${MAP_KEY}/${source}/${west},${south},${east},${north}/${days}/2024-01-15`

    console.log("Fetching from NASA FIRMS:", firmsUrl)

    const response = await fetch(firmsUrl)

    if (!response.ok) {
      throw new Error(`NASA FIRMS API error: ${response.status}`)
    }

    const csvText = await response.text()

    // Parse CSV data
    const lines = csvText.trim().split("\n")
    const headers = lines[0].split(",")

    const fires: FireData[] = lines.slice(1).map((line) => {
      const values = line.split(",")
      const fire: any = {}

      headers.forEach((header, index) => {
        const value = values[index]
        switch (header.trim()) {
          case "latitude":
          case "longitude":
          case "brightness":
          case "scan":
          case "track":
          case "confidence":
          case "bright_t31":
          case "frp":
            fire[header.trim()] = Number.parseFloat(value) || 0
            break
          default:
            fire[header.trim()] = value
        }
      })

      return fire as FireData
    })

    return json({
      fires,
      count: fires.length,
      bounds: { west, south, east, north },
      source,
      days: Number.parseInt(days),
    })
  } catch (error) {
    console.error("Error fetching fire data:", error)
    return json(
      { error: "Failed to fetch fire data", details: error instanceof Error ? error.message : "Unknown error" },
      { status: 500 },
    )
  }
}

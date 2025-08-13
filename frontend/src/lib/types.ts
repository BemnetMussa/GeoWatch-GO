export interface FireData {
  latitude: number
  longitude: number
  brightness: number
  scan: number
  track: number
  acq_date: string
  acq_time: string
  satellite: string
  confidence: number
  version: string
  bright_t31: number
  frp: number
  daynight: "D" | "N"
}

export interface MapBounds {
  west: number
  south: number
  east: number
  north: number
}

// <CHANGE> Added missing ChangeDetectionResult interface
export interface ChangeDetectionResult {
  newFires: FireData[]
  extinguishedFires: FireData[]
  growingFires: FireData[]
  diminishingFires: FireData[]
  summary: {
    totalNew: number
    totalExtinguished: number
    totalGrowing: number
    totalDiminishing: number
    changePercentage: number
  }
}

// <CHANGE> Added missing ChangeDetectionPeriod interface
export interface ChangeDetectionPeriod {
  label: string
  days: number
}

// <CHANGE> Added missing FireChange interface
export interface FireChange extends FireData {
  changeType: 'new' | 'extinguished' | 'growing' | 'diminishing'
  previousBrightness?: number
  brightnessChange?: number
}

export interface RealTimeSettings {
  enabled: boolean
  interval: number // minutes
  notifications: boolean
  autoRefreshOnFocus: boolean
  backgroundSync: boolean
}

export interface NotificationData {
  id: string
  type: "new_fire" | "fire_growth" | "fire_extinguished" | "system"
  title: string
  message: string
  timestamp: Date
  location?: { lat: number; lng: number }
  severity: "low" | "medium" | "high"
}

export interface ConnectionStatus {
  status: "connected" | "disconnected" | "error"
  lastSync?: Date
  nextSync?: Date
  updateCount: number
}

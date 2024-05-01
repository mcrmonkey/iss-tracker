#!/usr/bin/env python3

import json
import requests
from math import radians, cos, sin, asin, sqrt, atan2, pi


DEVICE_LAT = 53.480970
DEVICE_LON = -2.237150

## Head to geocode.maps.co and sign up for free account to obtain an API key
GEOCODE_API_KEY = ""

OCEAN_COUNTRY = "Ocean"
CITY_UNKNOWN = "Unknown City"

country = OCEAN_COUNTRY
city = CITY_UNKNOWN

def haversine(lat1, lon1, lat2, lon2):
    lon1, lat1, lon2, lat2 = map(radians, [lon1, lat1, lon2, lat2])
    dlon = lon2 - lon1
    dlat = lat2 - lat1
    a = sin(dlat/2)**2 + cos(lat1) * cos(lat2) * sin(dlon/2)**2
    c = 2 * asin(sqrt(a))
    return c * 3956

def direction_lookup(destination_x, destination_y, origin_x, origin_y):
    deltaX = destination_x - origin_x
    deltaY = destination_y - origin_y
    degrees_temp = atan2(deltaX, deltaY)/pi*180
    if degrees_temp < 0:
        degrees_final = 360 + degrees_temp
    else:
        degrees_final = degrees_temp
    compass_brackets = ["N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW", "N"]
    compass_lookup = round(degrees_final / 22.5)

    #compass_brackets = ["N", "NE", "E", "SE", "S", "SW", "W", "NW", "N"]
    #compass_lookup = round(degrees_final / 45)

    return compass_brackets[compass_lookup], degrees_final


response_doc = requests.get("http://api.open-notify.org/iss-now.json").json()
iss_lat = float(response_doc["iss_position"]["latitude"])
iss_lon = float(response_doc["iss_position"]["longitude"])
iss_distance = round(haversine(DEVICE_LAT, DEVICE_LON, iss_lat, iss_lon))


geo_doc = requests.get(f"https://geocode.maps.co/reverse?lat={iss_lat}&lon={iss_lon}&api_key={GEOCODE_API_KEY}").json()

try:
    country = geo_doc["address"]["country"][:13]
except Exception:
    pass

if country != OCEAN_COUNTRY:
    try:
        city = geo_doc["address"]["city"][:13]
    except Exception:
        try:
            city = geo_doc["address"]["suburb"][:13]
        except Exception:
            try:
                city = geo_doc["address"]["state"][:13]
            except Exception:
                pass





print("ISS distance from your location is: " + str(iss_distance) + " mi")
print("Country: " + country )
print("Nearest city: " + city )

direction = direction_lookup(iss_lon, iss_lat, DEVICE_LON,DEVICE_LAT)

print("Viewing direction: " + direction[0] + " (Bearing: " + str(direction[1]) + ")" )



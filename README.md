
# ISS Tracker

Tracks where the ISS is and tells you which direction you'll need to look to
see it.

Written in both Python and Go

Also tells you where on earth its currently flying over

Some of the math may be slightly out so verify stuff before relying on it
entirely

Based on some work at: github.com/simonprickett/pico-gfx-portal

### Geolocation 

Both versions make use of geocode.maps.co which now requires an API key
Sign up there for a free api key to use their service. There are some
limitations on its use so dont abuse it.


### Your location

Both versions default to Manchester, UK as the viewing location
Ensure to adjust this to your preferred viewing location for best results

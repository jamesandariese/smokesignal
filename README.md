# smokesignal

Tee results from a nagios check to a TSDB.  This adds no information to the
performance data but may add errors or important information to the extended
performance data.

Translates unknown to -1 unless you specify `-negative-unknown=false`.

## Usage

```
smokesignal
    [-influx-url "http://localhost:8086/nagios"]
    [-influx-json-tags '{"host": "sdf.org"}']
    [-negative-unknown=false]
    -measurement temperature
    command arg1 .. argN
```

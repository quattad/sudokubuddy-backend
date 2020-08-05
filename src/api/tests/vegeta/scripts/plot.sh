cd ../targets/

cat results.bin | vegeta plot > plot.html

start "Google Chrome" plot.html
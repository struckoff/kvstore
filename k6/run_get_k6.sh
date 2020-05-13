# K6_VUS=3
# K6_VUS=300
K6_VUS=50
# K6_ITERATIONS=259200
K6_ITERATIONS=25000
K6_KEYS=/data/datafile

./kvstoregeokeys -address "http://localhost:9191" -l 100 -n $K6_ITERATIONS > ./data/datafile

docker run -v `pwd`/data:/data --net=host -i loadimpact/k6 run -e K6_KEYS=$K6_KEYS -e K6_VUS=$K6_VUS -e K6_ITERATIONS=$K6_ITERATIONS --vus $K6_VUS --iterations=$K6_ITERATIONS - <"$1"

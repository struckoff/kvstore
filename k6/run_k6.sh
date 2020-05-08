K6_VUS=100
#K6_VUS=300
#K6_ITERATIONS=259200
K6_ITERATIONS=50000

docker run --net=host -i loadimpact/k6 run -e K6_VUS=$K6_VUS -e K6_ITERATIONS=$K6_ITERATIONS --vus $K6_VUS --iterations=$K6_ITERATIONS - <"$1"

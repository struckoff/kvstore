for port in $(seq 9190 9200); do
  echo "$port - "`(http :$port/nodes | jq -r ".[] | .ID"  | wc -l)`;
done

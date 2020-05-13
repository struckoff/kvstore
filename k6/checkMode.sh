for port in $(seq 9190 9200); do
  echo "$port - "`http :$port/config | jq .Mode`;
done

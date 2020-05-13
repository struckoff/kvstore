for port in $(seq 9190 9200); do
  echo "$port - "`http OPTIONS :$port/config/log/enable`;
done

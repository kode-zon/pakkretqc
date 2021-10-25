jsvalue=($(jq -r '.version' ../package.json))
echo "version = ${jsvalue[@]}"

docker run --rm -it \
   -p 0.0.0.0:8080:8888 \
   -v /app/pakkretqc/pakkretqc-master:/app/pakkretqc \
   --name pakkretqc --detach "pakkretqc:v${jsvalue[@]}" bash

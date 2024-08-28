
if [ ! -f ".air.toml" ]; then
  echo "Initializing air configuration..."
  air init
fi

air

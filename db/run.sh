DB_INIT_DIR="/docker-entrypoint-initdb.d"
INIT_SCRIPTS_DIR="init-scripts"

echo "Starting"

if [ ! -d "$DB_INIT_DIR" ]; then
  echo "$DB_INIT_DIR doesn't exist - creating"
  mkdir "$DB_INIT_DIR"
fi

echo "Copying init scripts to $DB_INIT_DIR"
cp "$INIT_SCRIPTS_DIR"/* "$DB_INIT_DIR/"

echo "Content of $DB_INIT_DIR"
echo "-----"
ls -al "$DB_INIT_DIR"
echo "-----"
echo "Done"

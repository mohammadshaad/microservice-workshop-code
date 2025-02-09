docker compose exec user-db psql -U postgres -d userdb -c "SELECT * FROM users;"

docker compose exec task-db mongosh taskdb --eval "db.tasks.find().pretty()"

![image](architecture.png)

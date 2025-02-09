![image](architecture.png)


# Build and deploy
docker-compose build
docker-compose up -d

# Check services
docker-compose ps
kubectl get pods
kubectl get services

# Check logs
kubectl logs deployment/user-interface
kubectl logs deployment/api-gateway


docker compose exec user-db psql -U postgres -d userdb -c "SELECT * FROM users;"

docker compose exec task-db mongosh taskdb --eval "db.tasks.find().pretty()"


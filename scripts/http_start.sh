#bin/bash
aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin 650142038379.dkr.ecr.ap-southeast-1.amazonaws.com
docker pull 650142038379.dkr.ecr.ap-southeast-1.amazonaws.com/salesman-service:http-$ENV_NAME-$IMAGE_VERSION
docker run --restart always --log-opt awslogs-stream=salesman-service-api -d -p 5014:8000 -p 5514:8001 --name salesman-service 650142038379.dkr.ecr.ap-southeast-1.amazonaws.com/salesman-service:http-$ENV_NAME-$IMAGE_VERSION main
cloudfront_id=""

if [ "$ENV_NAME" == "staging" ]; then
  cloudfront_id=$(aws secretsmanager get-secret-value --secret-id StgCloudFrontID --query SecretString --output text --region ap-southeast-1)
elif [ "$ENV_NAME" == "development" ]; then
  cloudfront_id=$(aws secretsmanager get-secret-value --secret-id DevCloudFrontID --query SecretString --output text --region ap-southeast-1)
elif [ "$ENV_NAME" == "production" ]; then
  cloudfront_id=$(aws secretsmanager get-secret-value --secret-id ProdCloudFrontID --query SecretString --output text --region ap-southeast-1)
fi

aws cloudfront create-invalidation --distribution-id "$cloudfront_id" --paths "/*"



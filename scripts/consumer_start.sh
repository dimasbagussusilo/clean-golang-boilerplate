#bin/bash
aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin 650142038379.dkr.ecr.ap-southeast-1.amazonaws.com
docker pull 650142038379.dkr.ecr.ap-southeast-1.amazonaws.com/salesman-service:consumer-$ENV_NAME-$IMAGE_VERSION
docker run --restart always --log-opt awslogs-stream=salesman-service-consumer -d --name consumer_$CONSUMER_NAME_ENV 650142038379.dkr.ecr.ap-southeast-1.amazonaws.com/salesman-service:consumer-$ENV_NAME-$IMAGE_VERSION -name=$CONSUMER_NAME_ENV consumer
cloudfront_id=""

# cp /home/ec2-user/salesman-service/supervisor.d/* /etc/supervisord.d/
# systemctl restart supervisord

# rm -rf /etc/awslogs/awslogs.conf
# cp /etc/awslogs/awslogs.conf.orig /etc/awslogs/awslogs.conf
# cp -fr /home/ec2-user/salesman-service/awslog/*.conf /etc/awslogs/config
# cat /etc/awslogs/config/*.conf >> /etc/awslogs/awslogs.conf
# systemctl restart awslogsd

if [ "$ENV_NAME" == "staging" ]; then
  cloudfront_id=$(aws secretsmanager get-secret-value --secret-id StgCloudFrontID --query SecretString --output text --region ap-southeast-1)
elif [ "$ENV_NAME" == "development" ]; then
  cloudfront_id=$(aws secretsmanager get-secret-value --secret-id DevCloudFrontID --query SecretString --output text --region ap-southeast-1)
elif [ "$ENV_NAME" == "production" ]; then
  cloudfront_id=$(aws secretsmanager get-secret-value --secret-id ProdCloudFrontID --query SecretString --output text --region ap-southeast-1)
fi

aws cloudfront create-invalidation --distribution-id "$cloudfront_id" --paths "/*"



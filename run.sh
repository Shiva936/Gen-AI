#!/bin/bash

# Function to display usage information
usage() {
  echo "Usage: $0 {start|stop|clean}"
  exit 1
}

# Function to start the services
start_services() {
  # Applying Kubernetes Application Deployment 
    echo "Building App Deployemnts..."
    kubectl apply -f myapp-deployment.yaml

    # Applying Kubernetes Horizontal Scaler  
    echo "Building Horizontal Scaler..."
    kubectl apply -f myapp-hpa.yaml
    
    # Build your service
    echo "Building Service..."
    docker-compose -f docker-compose.yaml build

    # Run Docker Compose
    echo "Running Docker Compose..."
    docker-compose -f docker-compose.yaml up -d

    echo "Deployment and Docker Compose setup complete!"
}

# Function to stop and clean the services
stop_services() {

    # Deletig Kubernetes Application Deployment 
    echo "Removing App Deployments..."
    kubectl delete -f myapp-deployment.yaml

    # Deleting Kubernetes Horizontal Scaler  
    echo "Removing Horizaontal Scaler..."
    kubectl delete -f myapp-hpa.yaml
    
    # Stop and remove Docker Compose services
    echo "Stoping Service..."
    docker-compose -f docker-compose.yaml down
    
    echo "Services stopped and cleaned up."
}

# Function to remove docekr image
remove_image(){
    # Removing the Docker image
    echo "Removing Docker Image."
    read -p "Do you want to proceed? (YES/NO) " yn
        case $yn in
            [Yy]* ) 
                docker rmi myapp:latest  
                echo "Docker Image Deleted."
            ;;
            [Nn]* ) return 1;;
            * ) echo "Please answer YES, NO, or CANCEL.";;
        esac
}

# Check for the correct number of arguments
if [ "$#" -ne 1 ]; then
  usage
fi

# Check the provided argument and call the appropriate function
case "$1" in
  start)
    start_services
    ;;
  stop)
    stop_services
    ;;
  clean)
    remove_image
    ;;
  *)
    usage
    ;;
esac

exit 0
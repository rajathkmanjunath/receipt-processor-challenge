# Receipt Processor Web Service

This is a solution created for the Fetch Rewards Backend Engineer assessment and is intended solely for the purpose of evaluation.

 This is a Go-based web service for processing receipts and calculating points based on specific rules.

## Getting Started

Follow these steps to set up and run the project in a Docker container.

### Prerequisites

- Docker: Ensure that Docker is installed on your system.

### Building the Docker Image

1. Clone or download this repository to your local machine.

2. Open a terminal and navigate to the project directory.

3. Build the Docker image using the provided Dockerfile:

   ```bash
   docker build -t receipt-processor .
   ```

### Running the Docker Container

    Run a Docker container from the built image:

   ```bash
    docker run -p 8080:8080 receipt-processor
   ```

This command starts the web service inside a Docker container, mapping port 8080 from the host to the container.

### API Endpoints

#### Process Receipts

    Endpoint: /receipts/process
    Method: POST
    Payload: Receipt JSON
    Response: JSON containing an ID for the receipt

#### Get Points

    Endpoint: /receipts/{id}/points
    Method: GET
    Response: JSON object containing the number of points awarded

### Testing

After setting up the service, you can test the endpoints as follows using one of the provided test receipts:
#### Adding a Receipt
1. To add a receipt, run the following `curl` command:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{
      "retailer": "Target",
      "purchaseDate": "2022-01-02",
      "purchaseTime": "13:13",
      "total": "1.25",
      "items": [
          {"shortDescription": "Pepsi - 12-oz", "price": "1.25"}
      ]
  }' http://localhost:8080/receipts/process
  ```
#### Getting Points for a Receipt
To get the points for a receipt, use the curl command with the receipt ID generated from the response of the previous command:
  ```bash
  curl http://localhost:8080/receipts/{reciept_id}/points
  ```
  Replace {receipt_id} with the actual ID generated from the response.
# Arvan Challenge - Data Processor
![image](https://github.com/raman-vhd/Data-Processor/assets/73130097/030be86b-5186-4429-bb36-93f892f2369c)
# Installation
```
docker-compose up
```
# Usage
- POST /api/data {id: DATAID, user_id: USERID}

# Technologies Used
- Redis: Employed for efficient rate limit management and duplicate checking, ensuring data integrity.
- MongoDB: Utilized as the primary database for storing and retrieving data.
- Kafka: Implemented to establish a data processing queue, ensuring seamless data flow and persistent storage.

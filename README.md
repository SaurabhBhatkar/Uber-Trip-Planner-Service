# Uber-Trip-Planner-Service
Trip Planner Part II(using UBER Price Estimate API &amp; UBER Sandbox environment for all API calls)

This project aims at finding the shortest route from source to destination for a round trip using Uber and Google maps API. 

User enters set of location address into the database. The co-ordinates for these locations are computed via Google maps API. 

The Users can perform CRUD operations like: 

Create / Post: Adding locations in the database.

GET: Return the information pertaining to a location w.r.t location ID via GO struct.

PUT: Update the information pertaining to a location w.r.t location ID with the values entered by the user.

Delete: Delete a location on the basis of location ID specified by the user

User specifies a starting point id and the list of intermediate locations.

The system computes the shortest route from the starting point and covers all the intermediate locations making a round trip reaching back to starting point. The Uber Price Estimate API is used to make this calculation. 

Apart from the best route, details like total uber cost, distance and total duration will also be sent back to the user. The status will be set to “planning” in this step.

User then starts the trip by requesting UBER for the first destination. This is done by UBER request API which requests a car from starting point to the next destination. The status will be changed to “request” in every requests.

The system is built in GO language and MongoDB is used for Data persistence.

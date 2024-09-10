To get started:
  run package "hub" .Main() method and provide your bridges ip. 
  Follow the instructions and copy out the access token provided.



TBD



next open topics:
-get websocket working >> review and test the code in hub package
-get basic device (light) working. python uses inheritance and kwargs serialization to roll out th response data from the get request into its light(device) class. maybe roll a custom util package which will do something similar > pass a response with kwards into a struct constructor and return the finished struct?

Demo app to test some of Elasticsearch's querying capabilities. You can set up
the required stack with the following repo: https://github.com/Omar-Khawaja/mysql-elk-docker-compose

After running the app, visit http://localhost:8080/home to test it out.

If you want to use curl to query ES in a way similar to how this app is doing,
you can run the following command:

`curl localhost:9200/poems/_search?pretty -H "Content-Type: application/json" -d @query.json`

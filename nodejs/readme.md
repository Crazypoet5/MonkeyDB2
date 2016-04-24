# Monkey-Node-Driver

## Usage

Like those shown in `test.js`:
```
var monkey = require("./monkey")
monkey.connect(function() {
    monkey.query("select * from t1", function(data) {
        console.log(data)
        monkey.close()
    })
})
```
You can run slowly in this version, and we will perfect it in future versions.

## Response

Response is a json string like:  
```
{"relation":[{"id":1,"name":"InsZVA","gpa":2},{"id":2,"name":"naloy","gpa":5}],"result":{"affectedRows":2,"usedTime":0},
"error":null}
```
realtion is the query relation, result is something around this query.

## Note

You should close `monkey` after your using, otherwise your node will never stop.
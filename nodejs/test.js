var monkey = require("./monkey")
monkey.connect(function() {
    monkey.query("select * from t1", function(data) {
        data = JSON.parse(data)
        console.log(data)
        monkey.close()
    })
})
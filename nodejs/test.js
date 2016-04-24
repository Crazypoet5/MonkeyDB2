var monkey = require("./monkey")
monkey.connect(function() {
    monkey.query("select * from t1", function(data) {
        console.log(data)
        monkey.close()
    })
})
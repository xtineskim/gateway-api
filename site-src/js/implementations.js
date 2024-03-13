document$.subscribe(function() {
  var tables = document.querySelectorAll("article table:not([class])")
  tables.forEach(function(table) {
    new Tablesort(table)
  })
})


//var conformance_path = "../../conformance/reports/v1.0.0/"
//var temp = "https://raw.githubusercontent.com/kubernetes-sigs/gateway-api/main/conformance/reports/v1.0.0/cilium.yaml"
//
//parseYaml(temp)
//function parseYaml(file) {
//  var request = new XMLHttpRequest();
//  request.open("GET", file, false);
//  request.send(null);
//  var returnValue = request.responseText;
//
//  test = returnValue
//
//  console.log(test[2])
//}

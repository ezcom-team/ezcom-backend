
pm.test("Response status code is 200", function () {
    pm.response.to.have.status(200);
});


pm.test("Response content type is text/xml", function () {
    pm.expect(pm.response.headers.get("Content-Type")).to.include("text/xml");
});


pm.test("Response body is not empty", function () {
    const responseData = xml2Json(pm.response.text());
    
    pm.expect(responseData).to.exist.and.to.not.be.empty;
});


pm.test("Response body is in a valid XML format", function () {
  const responseData = xml2Json(pm.response.text());

  pm.expect(responseData).to.not.be.null;
});

// Parse the response body to extract the productId and store it in a global variable
var jsonData = xml2Json(pm.response.text());
var productId = jsonData.productId; // Assuming the productId is directly accessible in the XML response
pm.globals.set("productId", productId);
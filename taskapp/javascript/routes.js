//SPDX-License-Identifier: Apache-2.0

var product = require('./controller.js');

module.exports = function(app){

  app.get('/get_product/:product', function(req, res){
    product.query(req, res);
  });
  app.post('/add_product', function(req, res){
    console.log(req.body);
    product.invoke(req, res);
  });

  app.post('/test', function(req, res){
    console.log("Req",req.body);
    res.send('Test');
  });
}

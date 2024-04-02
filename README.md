# OData Query processing
This is a library to process OData style queries especially $filter. It is not intended to be complete, but rather to cover the more reasonable aspects of the standard.

Examples of what is not covered:
* Filtering inside of $orderby directives. i.e. Microsoft example of 
```
http://host/service.svc/Orders?$orderby = ShipCountry ne 'France' desc
```
would instead need to be the following for this implementation
```
http://host/service.svc/Orders?$filter=ShipCountry ne 'France'&$orderby=ShipCountry desc
```
* Expand is not currently covered
* Some filter options including use of $it

If there are features you wish to add that don't compromise the simplicity of the code or majorly impact the speed please open a pull request.
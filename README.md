# viewer

## Setup

1. Clone the repository
2. Download the [Cars API](https://intra.hive.fi/api/file?gitpath=content/coding-fundamentals/coding-fundamentals-go/cars/viewer/_config/resources/api.zip&_gl=1*13z0oq1*_up*MQ..*_ga*NzcxOTE5ODk1LjE3ODI0NjIxMDQ.*_ga_YWXED0Y2T1*czE3ODI0NjIxMDMkbzEkZzAkdDE3ODI0NjIxMDMkajYwJGwwJGgw*_ga_ESP7H2ZXC7*czE3ODI0NjIxMDQkbzEkZzAkdDE3ODI0NjIxMDQkajYwJGwwJGgxMDM3NzI5ODY1)
3. go (v1.26 or later) installed
4. Node.js (v24.18.0 or later) installed
5. Npm (v11.16.0 or later) installed

6. In terminal navigate to your downloaded api folder. Launch your api server in terminal with commands ```make build``` and ```make run```
7. Open another terminal. Navigate to your cloned project folder and launch your web server with command ```go run .```
8. Web page can be accessed in a web browser with URL: ```http://localhost:8080/```



## Features
- Fetches car information from an API
- Home page lists all cars
- Full details for every specific car
- Advanced filtering with checkboxes for Manufacturer and Category
- Compare cars (max **2**) show the cars dtails side by side
- Recently viewed cars (max **5**) at bottom of the car page
- Recommendations based on the users recently viewed cars

## Usage
- Click a car to open a popup card
- From the popup card, open a full details page for that car
  - Chose to add or remove a car to compare
  - Recently viewed cars and Reommendations on to bottom of the page
- Filter cars on the home page (checkboxes) by:
  - Manufacturer
  - Category e.g., SUV, Sedan, Hatchback
- Choose compare to navigate to compare page

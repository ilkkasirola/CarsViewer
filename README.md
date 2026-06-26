# viewer
Website that showcases information about different car models and their specifications, manufacturers.

## Setup

1. Clone the repository
2. Download the [Cars API](https://intra.hive.fi/api/file?gitpath=content/coding-fundamentals/coding-fundamentals-go/cars/viewer/_config/resources/api.zip&_gl=1*13z0oq1*_up*MQ..*_ga*NzcxOTE5ODk1LjE3ODI0NjIxMDQ.*_ga_YWXED0Y2T1*czE3ODI0NjIxMDMkbzEkZzAkdDE3ODI0NjIxMDMkajYwJGwwJGgw*_ga_ESP7H2ZXC7*czE3ODI0NjIxMDQkbzEkZzAkdDE3ODI0NjIxMDQkajYwJGwwJGgxMDM3NzI5ODY1)
3. go (v1.25 or later) installed
4. Node.js installed
5. Npm installed

## Usage

1. In terminal navigate to your downloaded api folder. Launch your api server in terminal with commands ```make build``` and ```make run```
2. Open another terminal. Navigate to your cloned project folder and launch your web server with command ```go run .```
3. Web page can be accessed in a web browser with URL: ```http://localhost:8080/```

## Features

- Filtering
    - Filter cars on main page by choosing which models or categories to show.

- Comparison
    - Comparison page for two cars to compare side by side. 

- Recently viewed
    - Section showing up to 5 cars the user has recently viewed.

- Recommendations
    - Section showing up to 5 cars recommended for user by ranking attributes of recently viewed cars.
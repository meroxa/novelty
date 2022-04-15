require 'faker'
require 'pg'
require 'pry'

USER_ACTIVITIES = ["created new account", "logged in", "logged out", "updated their profile", "deleted their account"]
TIME_OF_DAY = ["morning", "afternoon", "evening", "night"]
LOCATIONS = [
    {country: "United States", city: "New York"},
    {country: "United States", city: "Los Angeles"},
    {country: "United States", city: "Chicago"},
    {country: "United States", city: "Houston"},
    {country: "United States", city: "Philadelphia"},
    {country: "United States", city: "Phoenix"},
    {country: "United States", city: "San Antonio"},
    {country: "United Kingdom", city: "London"},
    {country: "United Kingdom", city: "Birmingham"},
    {country: "United Kingdom", city: "Manchester"},
    {country: "United Kingdom", city: "Liverpool"},
    {country: "United Kingdom", city: "Bristol"}
]

class User
    attr_accessor :id, :username, :email, :first_name, :last_name, :country, :city

    def initialize
        @id = rand(1..100000)
        @username = Faker::Internet.user_name
        @email = Faker::Internet.email
        @first_name = Faker::Name.first_name
        @last_name = Faker::Name.last_name
        loc = LOCATIONS.sample
        @country = loc[:country]
        @city = loc[:city]
    end
end

USERS = []
for i in 1..100
    USERS[i] = User.new
end

def create_user_activity
    u = USERS.sample
    timestamp = Time.now.utc
    return [u.id, u.first_name, u.last_name, u.email, USER_ACTIVITIES[1..3].sample, timestamp, u.country, u.city]
end

pry
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
100.times do
    USERS.append(User.new)
end

def create_user_activity
    u = USERS.sample(1).first
    timestamp = Time.now.utc
    return [u.id, u.first_name, u.last_name, u.email, USER_ACTIVITIES[1..3].sample, timestamp, u.country, u.city]
end

def create_anomalous_location(country:, city:)
    u = USERS.sample(1).first
    timestamp = Time.now.utc
    activity = [u.id, u.first_name, u.last_name, u.email, USER_ACTIVITIES[0], timestamp, country, city]
    conn = PG.connect(ENV["DATABASE_URL"])
    conn.exec_params("INSERT INTO user_activity(user_id, first_name, last_name, email, activity, timestamp, country, city) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", activity)
end

def create_anomalous_timestamp(timestamp:)
    u = USERS.sample(1).first
    timestamp = Time.parse(timestamp).utc
    activity = [u.id, u.first_name, u.last_name, u.email, USER_ACTIVITIES[0], timestamp, u.country, u.city]
    conn = PG.connect(ENV["DATABASE_URL"])
    conn.exec_params("INSERT INTO user_activity(user_id, first_name, last_name, email, activity, timestamp, country, city) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", activity)
end

def generate_user_activity(count:)
    conn = PG.connect(ENV["DATABASE_URL"])
    count.times do
        conn.exec_params("INSERT INTO user_activity(user_id, first_name, last_name, email, activity, timestamp, country, city) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", create_user_activity)
    end
end

pry

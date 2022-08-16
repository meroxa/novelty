# Meroxa Data App for thatDot Novelty Detector

This Turbine Data App pulls user activity data from a Postgres database, transforms the data to make it categorical (optimizing it for use with the Novelty Detector), submits it to the Novelty Detector REST API and writes the response (scored data) back into Postgres.

It reads records from a table named: `user_activity` and will write to a table named `user_activity_enriched`. This new table will be automatically created if it doesn't already exist.

## Getting Started

### Prerequisites

* Postgres Database
* [Novelty Detector Server](https://www.thatdot.com/product/novelty-detector)
* [Meroxa Data Platform](https://meroxa.com) account

### Steps

1. Clone this repo: `git clone https://github.com/meroxa/novelty.git`
1. Change into the downloaded directory: `cd novelty`
1. Set environment variable for Novelty Server: `export NOVELTY_SERVER_URL=<url>`
1. Run locally to verify everything is working: `meroxa apps run`

At this point the Meroxa CLI will exercise the data app and produce an output showing the sample records (from the _fixtures_ directory enriched with Novelty data).

If everything's looking good, you can then deploy the Turbine app to the platform:
```bash
meroxa apps deploy
```

## Scripts

Included in this repo is a Ruby script (`data_gen.rb`) that uses the [Faker](https://github.com/faker-ruby/faker) gem to generate fake user data.

The script requires the environment variable `DATABASE_URL` to be set, so that it can write records into Postgres.

Launch REPL:
```bash
ruby scripts/data_gen.rb
```

Generate sample data:
```ruby
generate_user_activity(count: 1000) # creates 1000 user activity records
```

Generate anomalous record:
```ruby
create_anomalous_location(country: "Japan", city: "Tokyo")
```

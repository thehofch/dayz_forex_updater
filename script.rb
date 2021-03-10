require 'json'
require 'sqlite3'

config = JSON.parse(File.open('config.json').read)
puts "Loaded config file"

db_file_exists = File.exists?(config['db_file_path'])

db = SQLite3::Database.new(config['db_file_path'])

def create_db(config, db)
  puts "Creating database"
  db.execute <<-SQL
    create table rates (
      currency varchar(255),
      val int,
      ratetime varchar(255)
    );
  SQL
  puts "Database created successfully"
  puts "Creating default rates"
  db.execute("insert into rates (currency, val, ratetime) values (?, ?, ?)", [
    config['currency_name'],
    config['default_currency_rate'],
    Time.now.to_s
  ])
end

def load_last_rate(db)
  rows = []
  puts "Loading last rate"
  db.execute( "SELECT * FROM rates ORDER BY ROWID DESC Limit 1" ) do |row|
    rows << {
      currency: row[0],
      val: row[1],
      ratetime: row[2]
     }
  end
  puts "Last rate loaded"
  rows[0]
end

def calculate_new_rate(last_rate, config)
  last_val = last_rate[:val]

  low_threshhold = config['min_sell_value'].to_f / 100
  high_threshhold = config['max_sell_value'].to_f / 100

  new_sell_val = rand(low_threshhold..high_threshhold) * last_val

  low_threshhold_buy = 1.0 - (config['max_buy_discounter'].to_f / 100)
  high_threshhold_buy = 1.0 - (config['min_buy_discounter'].to_f / 100)

  new_buy_val = rand(low_threshhold_buy..high_threshhold_buy) * new_sell_val

  { new_sell_val: new_sell_val.round, new_buy_val: new_buy_val.round }
end

def update_db(new_rate, db, config)
  puts "Updating new rate"
  db.execute("insert into rates (currency, val, ratetime) values (?, ?, ?)", [
    config['currency_name'],
    new_rate[:new_sell_val],
    Time.now.to_s
  ])
  puts "New rate updated"
end

def update_trader_file(new_rate, db, config)
  puts "Updating trader file"
  File.open(config['forex_trader_file_path'], 'w') { |f|
    f.write("<Trader> #{config['trader_name']}\n<Category> Currency\n#{config['currency_name']}, *, *, #{new_rate[:new_sell_val]}, #{new_rate[:new_buy_val]}\n<FileEnd>")
    puts "Trader file updated"
  }
end

create_db(config, db) unless db_file_exists

last_rate = load_last_rate(db)
new_rate = calculate_new_rate(last_rate, config)

update_db(new_rate, db, config)
update_trader_file(new_rate, db, config)

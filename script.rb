# File.open('C:\Users\hofch\Desktop\Trader\z_Forex.txt', 'w') { |f|
#   f.write("<Trader> Forex Trader\n<Category> Currency\nBitcoin, *, *, 100000, 96000\n<FileEnd>")
# }

require 'json'
require 'sqlite3'

config = JSON.parse(File.open('config.json').read)
puts "Loaded config file"


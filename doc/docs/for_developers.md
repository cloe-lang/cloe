# For Developers

## Preparing environment

Install Go 1.8+ and `rake`, the task runner in Ruby.
Then, install Tisp's dependnecies.

```
rake install_deps
```

## Testing

- `rake unit_test` runs unit tests.
- `rake command_test` runs command tests.
- `rake test` runs both.

## Other utility tasks

- `rake format` formats all files.
- `rake lint` lints all files.
- `rake clean` cleans up the repository directory.

For more information, see [rakefile.rb](https://github.com/tisp-lang/tisp/blob/master/rakefile.rb).

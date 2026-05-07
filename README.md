# Mockdata

`mockdata` is a toy frontend for the [datastream](https://github.com/mshindle/datastream) library. It provides various examples of how to generate mock data and push it into different systems.

## Features

This project demonstrates several data generation scenarios:

- **Craps**: Simulates dice rolls for a craps game.
- **Factory**: Generates mock factory production data.
- **Mobile Logs**: Generates simulated log entries from mobile devices.
- **Performance**: Tools for performance testing and data throughput simulation.

## Usage

You can run the different simulations using the subcommands provided by the `mockdata` binary:

```bash
# General help
./mockdata --help

# Run specific simulations
./mockdata craps
./mockdata factory
./mockdata mobile-logs
./mockdata perf
```

## Configuration

`mockdata` supports configuration via environment variables (prefixed with `MD_`) and `.env` files.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

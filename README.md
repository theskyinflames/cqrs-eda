# CQRS - EDA
This repo contains a set of tools to implement CQRS/EDA services. This tooling is composed of:

* CQRS utils:
    * Command, command handler
    * Command middleware
    * Query, query handler
    * Query handler middleware
* EDA:
    * Events basic
    * Events listener
* Bus:
    * Sequential generic bus
    * Concurrent generic bus

## CQRS
[CQRS](https://learn.microsoft.com/en-us/azure/architecture/patterns/cqrs) is a pattern that allows isolating the operations that modify the domain state, called *Commands*, from those that don't, called *Queries*. As a result of a *Command* execution, one or more domain events will be published.

As commands as queries are handled by specialized handlers called *Command Handler* and *Query Handler* respectively.

In some documentation, CQRS uses a read model as a separate infrastructure that serves queries, but I don't see it this way. You can use  CQRS without that, and it only makes sense when you need to split R/W operations on your domain.

### C/Q handler middlewares
CommandHandler and QueryHandler middlewares are used to intercept the flux to and from the handler. They help inject handler dependencies and react to the handler return. There is a simple example of middleware that prints handler-returned errors. Another option is to wrap command handlers in a DB TX middleware in charge of starting a transaction, pass it to the command handler, and rollbacking or committing the tx depending on whether the command handler fails.

It is a pattern that allows the C/Q handler taking care only of what is its responsibility as application services

You will find the CQRS tooling in [pkg/cqrs](pkg/cqrs) directory.

## Events and EDA
[EDA](https://en.wikipedia.org/wiki/Event-driven_architecture) stands for *Event-Driven-Architecture* It's an architectural pattern that allows decoupling the command handler that executes the command, and hence, the one that changes the domain, from those that react to this change. These reacting command handlers can belong to the same service or not. This decoupling is achieved by domain events publishing.

You will find the Events tooling in [pkg/events](pkg/events) directory.

### Events listener
The events tooling includes an events listener implementation. It's in charge of listening to a specific event and dispatching it to an event handler. Usually, this event handler will map the event to a command and call a command handler to react to the domain change notified by the event.

Take into account that the tradeoff of EDA architectures is [eventual consistency](https://en.wikipedia.org/wiki/Eventual_consistency)

## Bus and Hexagonal Architecture
As CQRS architectures as EDA ones, use a bus:

* CQRS uses a command/query bus to dispatch the commands and queries from the entry point, usually an HTTP or RCP API. It allows the decoupling of the infra layer from the application layer where command and query handlers live.

* For EDA architectures, when event consumers live in the same bounded context (same service), they're usually dispatched to an Even bus, which is in charge of delegating them to the corresponding event handlers. Usually, these event handlers will map the event to a command and dispatch it to the command bus.

You will find the Bus tooling in [pkg/bus](pkg/bus) directory.

### Bus implementations
There are two bus implementations: a sequential and a concurrent. Use the first one if you don't have performance issues related to events dispatching.

## Examples
I've implemented some examples to help you to understand how to use this tooling:

* [command_bus](examples/concurrent_bus) Example of a command dispatched to a command bus using the sequential bus.

* [event_bus](examples/events_bus) Example of an event dispatched to a sequential event bus.

* [concurrent_event_bus](examples/concurrent_event_bus) Example of two events dispatched to a concurrent event bus

## Do you think this is useful? back me up
Thinking and building this tool has taken part of my time and effort. If you find it useful, and you think I deserve it, you can invite me a coffee :-)

 <a href="https://www.buymeacoffee.com/jaumearus" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" height="41" width="174"></a>

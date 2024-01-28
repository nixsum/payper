# payper

Payper is a tool used for forwarding a TCP or UNIX socket from one machine to another.
Using this you can set up a backend that is behind a NAT, with the idea that the backend reaches out to the frontend, which is usually not the case.

It would look something like this, the {} braces represent a machine, [] braces a process and () a protocol(s) along the arrows:

[Client]----(tcp)---->{[Frontend]---(tcp/unix)--->[payper listener]<---}-----(tcp)-------{NAT}-------{[payper connector]-----(tcp/unix)----->[backend]}


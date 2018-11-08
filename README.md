# huebridge-shell

Most of my media center has been configured so I can control it from
an Alexa skill

>> Alexa tell media to pause

The code works out what is currently playing (BluRay, TiVo, iTunes,
DVD Player app, etc etc) and sends the necessary commands.

However, skills can't be used from routines.  I can't create a routine
that will do the equivalent of

>> Alexa lights on

>> Alexa tell media to play music

This doesn't seem to be possible.

But Smart Home devices _can_ be called.

There are a number of programs out there that provide a form of Philips
Hue smart light emulation.  They're good enough that you can use them
from Alexa, for example.

So I decided to take that idea, but rather than model the lights internally,
I would use an external program to do all the hard work.  This would allow
me to emulate on/off status of a light by the state of the receiver; for
example, the BluRay light would be on if the BluRay input was selected.
In this way the state of the emulated lights would reflect the real state
of the media center.

Some simple actions could also be added; eg turning on a "play music"
light could perform a set of action (turn on the receiver, start iTunes
playing).  Speciality lights ("play Christmas music") can be added as
well.

And so this code...  a bridge between the network communications needed
to pretend to be a Hue Bridge and a program that does the backend work.

The bridge code basically passes Hue API requests through to a backend
process; the backend process is responsible to reporting state.  For
efficiently the bridge caches state and will return that to any Hue API
calls.  The backend process can refresh the cache at any time (e.g. if
it detects an external event... a volume change on the receiver could
update the "brightness" of the associated light).

I also provide a dummy emulated media center showing how this works.

## Thanks to

I learned a lot by reading the code at https://github.com/pborges/huejack
and https://github.com/mdempsky/huejack - in particular how the UPNP
stuff works, and some nice ideas on how to present the XML and JSON
structures.  I'm probably not using idiomatic GO coding styles, but
it works for me!

## Communication between bridge and program

The communication patterns are asynchronous.  The bridge can send
commands to the program and the program can send status updates back
to the bridge at any time.  It's a very simple protocol.

### Commands:

There's only one command:

LIGHT#name#on/off#brightness  -- Set the light named "name"
to the on/off status and the brightness (0-255).  The name must not
contain a # since that is used as the separator.
If either the on/off or brightness values are - then that value should
be left unchanged.

example:  LIGHT#Play Music#on#-

Case of the light name should be insenstive, but that's really up to
the program; the bridge will always report what it's learned, but will
pass through whatever the API caller presents.

### Status updates

LIGHT#name#on/off#brightness -- This is the current state of
the light.  

LIST#name1#name2#name3 ...  -- This can be used to tell the
bridge of the complete set of lights being controlled.  For example, if
you have an environment where devices may join and leave (a mobile
phone, perhaps?) then this can be used as a way of refreshing the bridge's
knowledge and to stop telling clients about lights that no longer exist.

The names, here, _are_ case sensitive.  If the client reports FOO in
the "LIGHT" output, but Foo in the "LIST" output then these will be
considered two different lights.

### Thoughts and commentary

Light names must not contain the # character.  This will break things.

The LIGHT status update should probably be sent after a command for
that light.  It may also be sent periodically so the bridge knows the
state of the world or if the program detects a change to the real world
state.

I'm not sure if the response to a change is correctly formatted.  Or
possibly it's just the Alexa app ignoring the response and believing
what it sent was applied...  but eventually it polls and updates with
the correct data.

### Security

There isn't really any security on this; the API will accept any user ID,
so anyone who can reach the endpoint can control the virtual lights.  But
since the goal is to make these controllable by Alexa, those people could
just use their voice :-)

### License

Apparently some of the history of this code goes all the way back to
https://github.com/armzilla/amazon-echo-ha-bridge which was under
Apache 2.0, so that's what I'm licensing this under.

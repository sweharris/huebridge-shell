This is a real setup I use for my media center.

This exact code is unlikely to be of direct use to people, but may provide
some hints as to how the huebridge can be used for complicated setups.

It has some dependencies I don't provide here, like the "Media_API"
command, the "tv" command, "itunes" command, and the "denon" command.
These are ways I can talk to parts of my media setup.

I have a Denon receiver.  The "denon" command can be used to change the
state of the receiver, and also to report on the current state (volume,
current selected input, etc).

We create a number of fake lights:

      "MediaPower" - This indicates if the receiver is on or off.
                     If a request comes in to turn the receiver on/off then
                     we do this.
                     Brightness of this fake light is mapped to volume.
                     Turning this off will use the Media_API to shut down
                     my whole media player

           "Music" - turning this light on will turn the system on, and
                     switch input to Mac and play the current iTunes queue.

  "ChristmasMusic" - turning this light on will turn the system on,
                     and start playing Christmas music.

     "RandomMusic" - turning this light on will turn the system on,
                     and tell itunes to play music from my collection
                     (close to the old "iTunes DJ" mode).

              "TV" - turning this light on/off will turn the TV on/off.

            "TiVo" - turning this light on will turn the system on and make
                     sure the input is set to TiVo

          "BluRay" - turning this light on will turn the system on and make
                     sure the input is set to BluRay

             "Mac" - turning this light on will turn the system on and make
                     sure the input is set to Mac

Turning off any light except for MediaPower or TV will pause itunes (if it's
playing).

The state of the Mac/BluRay/TiVo lights change according to the selected
input on the Receiver.  The state of the TV is determined by if the TV is
on or not.  The state of the receiver is based on the on/off state, and
the brightness is determined by the volume.

So what this all means is that I can do things with the normal remote
control, and the state of the system is "adequately" reflected in the
state of the fake lights.

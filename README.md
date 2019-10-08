
brys:
: Cornish for intention, mind, opinion, thought, psyche
: Welsh for speed, haste, urgency

This is a place to work out the creation of a wiki in golang that
meets two critical requirements:

* Teaches me a bit more go (and the developer context).
* Results in a wiki I actually like.

A significant portion of my long development career has been
associated with creating wikis and wiki-like tools for information
management. Modern wikis fail to be good because they are
frequently focused on publishing; a dumping ground for stuff that
needs to be remembered. That's not where wikis shine. Wikis are at
their best when they identify gaps in our mental namespace and help
them to be filled. "Our" is an individual or small group. The
outcome of an effective wiki is shared language, shared
understanding and shared goals.

WikiLinks provide the mechanism for identifying the gaps.

# Desired Features

_In no particular order._

* [ ] densely documented on the wiki way in code, commit and wiki itself
* [ ] markdown with WikiWords
* [ ] one go process per wiki (but with interwiki links?)
* [ ] storage at first in plain text (then in git?)
* [ ] one namespace, no nesting
* [ ] no tags
* [ ] WikiWord autocomplete in editor, but other than that, very little JavaScript
* [ ] no built in auth
* [ ] RecentChanges
* [ ] NeedAttention (stuff that hasn't been edited in a while)
* [ ] MissingLinks (the gaps)
* [ ] Ambient info
    * [ ] page heat
    * [ ] backlinks

# Things to Maybe do Later

* [ ] templated layout and style
* [ ] purple
* [ ] sessions/identity (for log authorship)

# Desired Development

* [ ] Test oriented/driven development
* [ ] Easy to deploy in multiple ways
* [ ] Fast feedback and iteration

# License

Chosen the GPL 3 because it most accurately reflects my politics
and attitudes with regard to software: If you want to use it, you
are a participant in its creation.

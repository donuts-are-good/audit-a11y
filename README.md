# 🌐 audit-a11y

audit-a11y is a command line tool for auditing web pages for accessibility concerns.

## 💡 why?

recently on shhhbb.com bbs there was a person who related their story with impending disability that was already affecting their accessibility capabilities. this program was an idea that came from that conversation.

## 🚀 usage

to use audit-a11y you can clone it from github and build it yourself, or download a release and run it right away like this:

```
./audit-a11y example.com
```

```
./audit-ally https://github.com
Accessibility report for URL: https://github.com
- ul element without a corresponding label element. (x23)
- Image element without an alt attribute. Line:  (x1)
- input element without a corresponding label element. (x5)
- Button element without a corresponding label element. (x1)
- Input element with type 'hidden' without a corresponding label element. (x3)
- Heading element without a corresponding title element. (x62)
- Title element not inside a head element. (x26)
- Form element without a submit button. (x2)
- Form element without an enctype attribute. (x2)
- Form element without proper validation attributes. (x2)
- Anchor element without an href attribute. Line:  (x1)
Total issues: 128
Total lines fetched: 2346
```

for each infraction, the infraction itself and a multiplier denoting the number of occurrences will be included. 

## 🤝 contributing

if you want to help out this project, great, i'd love to hear your ideas. also, thank you :)

some todo items that make sense to add are:

- printing the offending html lines under each infraction category
- count the lines seen vs the lines scanned
- include an MDN link for each infraction type for more information on fixing it
- spider mode, integrate `donuts-are-good/araknnid` to crawl everything on a specific domain and prepare a report for each page found in that domain.
  - requires using araknnid with `--depth` to narrow the scope, may need to whitelist rather than blacklist
- add more checks
- add unit tests for existing checks

## 📄 license
MIT License, 2023, donuts-are-good https://github.com/donuts-are-good

if you don't know what it means, don't sweat it
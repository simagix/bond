# Chen's Bond - MongoDB Sharded Cluster Analysis Tool
Chen's Bond is a MongoDB analysis tool designed for the purpose of database management and optimization. It primarily operates by querying the 'config' database, providing users with valuable insights into their MongoDB deployment. The tool is particularly focused on improving the performance and efficiency of MongoDB clusters.

## Scripts
**Intro:**
1.
Say hello to "Chen's Bond," your MongoDB cluster buddy! It's like having a best friend for your sharded clusters, here to spill the beans on what's happening in your MongoDB world. Think of it as the secret sauce that brings all your shards together, giving you the inside scoop on your database setup. With Bond, you're not just keeping an eye on shards; you're simplifying sharded cluster analysis!

**Installation:**
2.
Now, let's talk about setting it up; it's as easy as ordering a pizza on a lazy Sunday. To get started, mosey on over to the Bond GitHub repository and grab that source code. You'll need Git for this, so make sure Git is on your computer. Once you've got the source code, open your terminal or command prompt, go to the downloaded directory, and fire up the "build dot sh" script.
3.
This script, tested on Linux, macOS, and Windows, is like magic, but you'll also need a Go compiler (version 1.19 or higher) for it to work. If you're on Windows, just use Git Bash and make sure MinGW and the Go compiler are on board before building Bond. The compilation process is faster than a cheetah on roller skates, and voila, Bond is ready to make your sharded cluster analysis a breeze!

**Connecting with Bond:**
4-5.
Bond offers versatile connection options tailored to your environment. Whether you're running a sharded cluster through a mongos instance, or operating in an isolated setting, Bond has you covered. It seamlessly connects to your sharded cluster via a mongos instance for instant insights. And for those in isolated environments, Bond allows you to restore to a standalone server, using a backup from your config database. Just remember to include a backup of the mongo-s collection for a smooth transition. You'll even find a handy script for backup and restoration of a config database example, making the process as straightforward as it gets.

**Summary and Alerts:**
6-8.
Time to dive into the good stuff. Bond starts by taking a peek at your cluster's configurations and metadata, like a detective on a case. It's also the Sherlock Holmes, of mongo-s versions, sniffing out any version mismatches, and health states. Plus, it keeps a watchful eye on your shard setup, making sure the "maxSize" parameter is playing nice. It's like a bouncer at the club, ensuring everyone's on the guest list.
9.
It even gives the side-eye to the "ActionLog" collection to confirm it's the capped collection it claims to be. Similarly, it gives a thumbs-up to the "ChangeLog" collection, checking that it's in the capped collection club. To top it off, Bond cross-references your cluster's version with MongoDB's recommended upgrade list, offering you a backstage pass to all the juicy details.
10.
Moreover, it scans through the shards, databases, and collections metadata to meticulously uncover chunk distributions, and other essential information. Bond's dedication to comprehensive insights doesn't stop there; it seamlessly links this chunk distribution information to dynamic charts, providing you with a visually engaging and intuitive way, to understand your MongoDB cluster's data distribution trends.
11.  
Now, let's chat more about the Alerts, the part where Bond plays protector. Any issues or red flags it discovers during its detective work are showcased here. For example, if the cluster mongo version is in the recommended upgrade list, Bond would flag a warning sign, and list links to the juicy detail. It's like your friendly neighborhood superhero, swooping in to save the day by alerting you to potential problems that might need your attention.

**Bug Tracking:**
12.
Ah, bugs! Even the fanciest software gets 'em, and MongoDB isn't immune. That's where Bond comes to the rescue, like a bug-swatting ninja. It keeps an eye on known issues and recommends upgrades to squash those pesky critters. But here's the fun part, you can join the bug-hunting party, by updating the "tickets dot json" file, with any known problems you've encountered. It's like adding your own bugs to the mix, but in a good way.
13.
Inside that file, you'll find not just bug references, but also nifty links for juicy bug descriptions. So, whether you're squashing bugs, or adding some to the collection, Bond is your trusty side arm in the wild west of sharded cluster maintenance.

**RESTful API:**
14.
Bond offers a comprehensive and user-friendly feature: a RESTful API that seamlessly provides JSON data for integration with your applications. This API empowers you to effortlessly access and, incorporate essential information from Bond, into your own software solutions, ensuring a harmonious and data-driven sharded cluster analysis experience. It's a bridge that, connects Bond with your applications, opening up a world of possibilities for enhancing the functionality, and insights of your MongoDB ecosystem.

**Charting Capabilities:**
15-16.
Last but not least, let's talk more about charts!  Visuals are where the party's at,  and Bond knows it.  While it currently offers some cool charting features about chunk migration and splits, we're here to spill the beans, the charting feature is about to get a makeover that'll make Cinderella jealous.  Get ready for even fancier charts that will provide deeper insights into your cluster's performance, health, and trends. 
17-18.
Oh,  I almost forgot to mention that you can access charts conveniently through RESTful API calls.  For instance,  you can automate the process of saving charts in HTML format.  As an example,  the URL on the screen will generate a pie chart.  And here's another one illustrating the distribution of a collection among shards.  It's all about making your cluster analysis as seamless and informative as possible. 

**Thank You:**
19.
Thanks for watching, and I hope you will enjoy using Chen's Bond to simplify your sharded cluster analysis. It's your trusty companion in the world of MongoDB, always here to help you make sense of it all. Happy sharding!

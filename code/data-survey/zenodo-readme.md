# Raw Data Set: Finding Unsafe Go Code in the Wild

This set contains the raw project code and dependencies data that we used for analysis for our paper
"Uncovering the Hidden Dangers: Finding Unsafe Go Code in the Wild".


## File contents

 - `projects.tar.gz`: contains the raw project code at the exact revision that we used for analysis.
   Uncompressed, this is 500 projects and makes up about 18 GiB.
 - `modules.tar.gz`: contains the dependencies. This needs to be unpacked into `$GOPATH/pkg/mod` and
   is about 55 GiB uncompressed.


## Further information

Further documentation as well as our analysis scripts can be obtained from our GitHub repository at
[https://github.com/stg-tud/unsafe_go_study_results](https://github.com/stg-tud/unsafe_go_study_results).


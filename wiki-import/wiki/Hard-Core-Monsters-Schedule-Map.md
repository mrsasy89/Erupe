To use the below table you will need to read the mhfsch.bin and produce the file in index order. The tables may only align using End of Life bin. 
```c
//------------------------------------------------
//--- 010 Editor v12.0.1 Binary Template
//
//      File: 
//   Authors: 
//   Version: 
//   Purpose: 
//  Category: 
// File Mask: 
//  ID Bytes: 
//   History: 
//------------------------------------------------
local int i;
Printf("|index| Year|Month|Day|Hour| \n");
Printf("|--| --|--|--|--| \n");


//adjust the number until it fits
local int value = 37;
for(i = 0; i <value ; i++){

struct{
uint32 year;
uint32 month;
uint32 day;
uint32 hour;
Printf("|%d| %d| %d| %d| %d|\n",i,year,month,day,hour);
        FSkip(16); 
} schedule;
};
```

# GHC
Derived from is_GHC_check
|Name                 |Monster ID |mhfsch index |
| ------------------- | --------- | ----------- |
|Berukyurosu          |85         |19           |
|Abiorugu             |104        |26           |
|Kuarusepusu          |105        |24           |
|Black Gravios        |47         |20           |
|Lunastra             |64         |23           |
|Akura Vashimu        |83         |18           |
|Cephadrome           |8          |16           |
|Plesioth             |21         |17           |
|Gold Rathian         |42         |21           |
|White Monoblos       |44         |27           |

# HC Monsters
Derived from is_em_special
|Name                 |Monster ID |mhfsch index |
| ------------------- | --------- | ----------- |
|Cephadrome           |8          |16           |
|Gravios              |17         |11           |
|Pink Rathian         |37         |4            |
|Blue Yian Kut-Ku     |38         |6            |
|Gold Rathian         |42         |21           |
|Black Diablos        |43         |10           |
|White Monoblos       |44         |27           |
|Red Khezu            |45         |3            |
|Green Plesioth       |46         |25           |
|Black Gravios        |47         |20           |
|Kushala Daora        |54         |1            |
|Rusted Kushala Daora |60         |12           |
|Lunastra             |64         |23           |
|Bulldrome            |68         |7            |
|Akantor              |77         |8            |
|Bright Hypnoc        |78         |2            |
|Lavasioth Subspecies |79         |29           |
|White Hypnoc         |82         |28           |
|Akura Vashimu        |83         |18           |
|Berukyurosu          |85         |14           |
|Nono Orugaron        |92         |0            |
|Doragyurosu          |95         |9            |
|Gurenzeburu          |96         |13           |
|Rukodiora            |99         |30           |
|Kokomoa              |102        |15           |
|Abiorugu             |104        |26           |
|Kuarusepusu          |105        |24           |
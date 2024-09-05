package packed

import "github.com/gogf/gf/v2/os/gres"

func init() {
	if err := gres.Add("H4sIAAAAAAAC/+x9CTSV69v+Y4rMJ1ODlGQetimSeTqxzUQUITa2tnlKkgYaJBoUwkmFUIQITU6ReQoRZW5QEqLI+F+nc9J+sbf97r07/+/7fb+z1sk6Z7Xu67nu5x2e974v92WsR0XNDugAHegw5rcCWP9wg5XAF+XmibH3RYm7ebh7ICSc0b72vr7e6L1+vigf8x00gKLkaJtTo2lFZY0eolZM4lmtbo0ewkDfKN9vnJYagLk5Yz1auiZjqW3cAABOAABuuNVLwqGd3T28UfNQshF6zMckGbXGfdg5v/IfoGI2ppF8zHHqeun7d185Gtw7dTySt5u1Bgs/lm/g3OJ0U2PQXC2lXdOJ1ur5yer3m46oG7zaUl6W5Bzv9ltZcpwUo37BiMo9c3Xrc0ffsm32OOKH6Fgj7GiTojAXpifC9Z6t8VBeZDMm8PlE/Oi3FT/YDOgNPAkBAGTiZcO5iI2B/T6UExrzk0tVjV6dsIGueYWOfl0L7W/3KhXiDhS1HQY/gJ6df7OL+3ti8AFxLQIy1VbXMtCWMNCaRzK18KwPEEXWGJjT6daKV+iK6Zp7ViNMh2vrq1mQ4hXveunFavXrLOgkdMVEcky7G80rhISFhe42thrVG1Ghm6oRyFaEjr6ReAuiWl+cNmds1QU96dHRj+qYevbdmpqabA56et55u2LXxcXGxccByh8ErrRX8fkAAFKwCADQvYDAb4sI2Ht6fl966dE2p7/iYP/tn3EWJoJnqTgIe0+0+I//O58PXEEXZ1dw2aCL95X+8zjLqYARRlNjY+Nsc+NnQnU6tZUmQkJGz15/9hnzY2L29PWn0RNOSU7lv/Q2WsirqHDFxMT2smPJ6r2naOLu8acX0Sic4k9HWkvM1nVmqW/dXCQ9et4HaWbKMTY29mVsTGbVhTuPQ/e4y238keh7Tc1d7gCAk3i5CC3PBe+lI1RVjzSQMIFcQp6MwzvoDMVFqvUN3/T29lLO3/O3zDm+yQIAJPDu2abll2TviV523wAxXO090QgXFAbjQcRlIQUj/N9/Sjj/xMmLVNcrUWMNq4+zHaHjWZVu9ZlKSWfVlga5jc9vRlRaHqv85K8YfalxTuIUzzTFNknG0ZH9Q2uOvSzmkV4RNkC9O2eVZQrVS2a+cYuigsaewVDagA372dmr1354dfYqJttdLq3qwIUjsymrA91VA0Unztp7as302xhvZn+aeUdk+5lT1yY+PHCteof6cu/RbnOuks51X2fOVPXTWsTxyCjGavKYIZMFEWkbBMfOXy8Q6Unc/Tju2hqBTJ6RC2c5ugKObFS03vAZc0mpUkJ7lePcukGus6ypfrwpDME0PzZeooW7pw0AkEWBb39E4STQX4qILZKFh7B4lzC1OYbHpBipu4sO0VeziZ2xq2Hztbd3uPnNK+FkbGhp+i7ezvUmAsmRgh0drYFBzhxKyoB9wG1r5+PIvP17bk9lWjWkCrvQ3yg/uyKemr+QdliPS0knc9BKxWvSMr6P66Pv3Oa6/R8/KphWFDNpOCeeYw55XOoTdDq5N9QnbqRQ6wCawWxN08g99U9hkcMztpflN4sKnLOu48ndzR8ts/v2KvSk0IYDwd9W/kj7ZEPRpwwAQD/e24Jv+aS42DvsIyLdCMIiIxw83J3QzhKB9m6YeZTNtTmuJZLsYT0dRl46LVXXCjfbCa3fRrNtjIN9YNfDHY0PK3IdioNN7gSXKEx3rLliK7/j3ARl5Qhw3t8uFC5QxREtL3E+fupTxJVdfmsymaS5YtJiVt26p3a6I5bH4GG8iN6sy6Nv4z3e6/fc1uR6K/z5BcWeVFsRgz+D3/RvT4g4pjf7ecWdjd57AlZrpyRLjm2dMbrJ0BXm1IkOPERxzHcsYZdZ632DE14ZIRQ/cr16YOLELQDAM7y5JuDVgXb3RXm722OIeMCJEh4d4eDmSMSOSsND+Otf7NvH96zRvleSrE/rG2+s1mLx9xRvv857QsDgpWbA7rvJrdL55btMZhw2e1DaFwflPDhr/XXz29qm+A2IngkWpsf9gsKZRucYZE3TpI5r1GvrCkfnvCmcHizvzz/l1OLUciPqgNNY8pl7rnQ78hQunc/avbfFaczs4edPqPG+0b7xE6pTRvdP3CmLlhWYenTWiJkBvY7eS+3GTtnIut8G+BS4LmJSlHhvv1LMo3C3Zxocuj6+OpVLuJRHmQlUXWfKuKCh5CQr++C3Ty8alQ5X0RT9iWSKm9CyLrC9JMksibDHlH1mjPGenBoc3MjrvLcwPSCk5WrY8dgjsXwF2RonGUvEY2d31rq1i3yxVZBIN9nlt/uRlXNh+tMvJwrZe12LHE/eGwrwv8ytzGt2R8E2zGO1wPaDSQl7K5offfNzLmlutmp+O6msPDPbkGfzPGP04ONtYvZ8jiIMtHQturx1iezObVVz81dhdSrjVRsKAOQo8V0nEnB20cPdx9eHiEtlK2yQf35gXzAiOvp1hvpG5jp19aISwm9/HiX+5GKs4QAAsOC9IaThrcHX2wODQXkTcedtJQqI6KOGOglwSxw96h6vpJRiDesJytjPVJj+54XsTfvzjnitlz/m8AdGQOceV1DD3GxJjMzXYok3EXo5K/ad5XmVucvvQsO6LKrtdfFKHCdGV6pvsww+7Pds5Z0ivkcDHda8TMfWe+xDIUPOpcaIsVhxaHfEqcxfpOnGOscxAABmCnwctUnmaOuOCoAesR4jS9TYtVuDbb2ENw0+YfwGuM5rMeVVpjqJ1r+2OnYhzv6gOjI2JKa/jKXhftGrpz08Dp5TrbwOxhPqnpyC58xXZK+THOmP3i5ItUdi//D4OUF1r7SrnS8ZUadTla+KDJitmOSSon9VPGo1LNBy4JoXvYqIf17bHvHyuwe37aiv7chBHteU5C5bd1HY0sLqxLfUu3mCjXqs7J3osRbkR++DLglBNqueDXw8IPu57c39mXoTC9Ws528sXdi03/V7zsTvWcU/Sf8jhX80mr0uBwCcxptCJOkp9JeyXXS96J9RNxqVZDxe/+X2MFfwCae1p8pU1RxOUlGhV308qeF5J1y20rogQ3aw+DonM7dYdcvoIyOv2WqaJ2xPYynasxw5b101VRG8qsu+QyTPe61KRH0ik4KsnKz/g3ym+5EYY3BxMqSt+FWxymDnlaGv92/6pVjubF93NUtLUHaDv8kLr+tqVLKany2DrN18VtVJrSp2k1Wp0x9EPbF92Db7cYt0otcW95qj3P23cpw+zD2a3LJqVCluy54k+rdq744cCOx+5Ya6fDpfwdjW+0MldxFHeLqFHk9UXN871ZSre88fYbPuUfyRYobub596//qEpCDXK9fR3uPXnFmxEb4XSPahUJ7zUPP33V/RaJZ9lInDwMN4OKMdiOAkBxfjX2Xl5uGIwvxiVt8xSGYlCRvRkZhLUIEIGJK5ycIGRbn7on0DieCnTCQUyRzhHMQ87R32oRx/8UHsb5B/fix1EAPP0oWtf5YjR1e4Pmdb9hiGgLECH5S3P9oB9Yt5/oMCcwMXYwosj+lmj3bHzmREnYFriSSrWk/+O3YNX17xjCOdSXmp7L3UlFbxp8REN7Ww1G1dHTHE85HNcW5if9jHxAds5zQo6XJenws7msHx2QFtWByj2cFTHH9AW4QmPdTIP4exlFlzr7zAWk8Za0EKkfdUXak3cxSu1T/a3tU3yjBK9fBNIKLTwu3jZpfDh5D+XWa7xIzbXsQd15B+J5xusCUkcs6se8MfF1uyn4j+4F5y7eveKwCAJlK/Zd3s3dFOKB9fIk7UEoRH/6d+QMQlsw02yJK1CtkIE3pKKUatOZ9wXtqGQ0jDHg5KpgHki0vJu28jOAIeOe2v0TuN5qV2ybSPO8pxabtw2EqXDubBKbNyvT8DBqcGEwcthVTvuXFyTRoJfL5bdf6egQ+dZUdIjMlLBWo+Q75Nk/N9Btt05kgnAEAI2XLniPLEeAT+mk+eBSCIfX4+vh5u6AMoIuCUSYBD7LX3IeaJYkAq5j//1w3l7gu9YnyjHhqWSzKWjis9fagWds9Jc69jf6m7KYNYAsvTFc0F2ynCdzXPGZjUjTet3xLY626qR53vFLDBZ/+F3N+VPJvuoSR4pBUzX7NTGQv5UI+7XvBxR9ZRGB9wyu9xsJvdYNknW1NnwH9Ou4ZC0sJwz5SR2IUbreVU785Tf6zvar7afvYPIM/65x4Xoee9m8OmeOWCh7S0BIRDig87HiiKDcmNCXiekpEs3rtHFX2H9dF848bTMP7mbQCAM95PDmOSU/bjP+190R7u0KyJVeUwPZFkpWmYTR2mlEKeeECp/NBLpb96xOa1z6aw98I635Jj5FWeJaYKN3ec2HLot9Pv469YDvGgCxwO2b+2EvHT7fHKVHIRWe8ehBjebGaZnekv13/S2f2m/uNnLJzTwpuK599rT+kqy/66zzwAid9Xy5D95z0EpWl9RoqeUp1ReyphE5dlg2HjgSM0T1VqDjl6Oq+Stq9HKJ6xG58YKjG7jShNawmLEg127QrdcPI6p7Hmh9R0jRqbdaG/yV3SGrVhk+V+4+t+ri8UdW14UKOmhueSBYbrUGdrUXLZuqJZQypF5K6K5g+zXUke87QTXRLO+gMAruG9FdVJoe3hj/LG2Af6EPEI0CUHLsIR5Y/CeHgS8TiwJCf+Py8SN3tP6OZjogz0StRYjw9949AMu3aiLPWZxYrwozwDBg+M7lzI9v1a0qF7tTsqpmwH0/6E8WMBPOFfE83uPubjvtlr91B9P/Oj/A7Rp09yzvPyiwElhd+72T9rtGnzvtGW2bz+ae/5SYsqNgX1CTEAsrrDor9GB/pOabkeQ9C2RvAwcx0OzerzLQKhUWJ0knzn2qIuPw/d85vsF6Xw47s/aM/JZqomnP4UrUptRhf+POPRfVsPnrdG7aUDrvnCMW9yb2VnxQzEv8wKjm7ABFzljduOjOu61TcY2BfQNvxoqKznNw/5B9MbXrE7qVEkJOmr/XyqIGZUAgcAAJF4nyp2ZM38z4ey+F/3H8p7we1XK5VRosZI3XDISIpL8kpXRt2aRoYguk2XLPaH8RawMGPC01Z52P3e9Enp+BdWl7cTinMvm/ayvYoNClG4fvmOb5MpWi3N7M/rUSn9Rftyzrw+gGLL2hB7oKP9evPR8lhKsSsSdNbij+hjB/oCNOT+eBX0iOpHLg4UXxbxAwDcxnsVWpM1F/ietjUGrqFqrNtHghOvIexO3W/YSWdFk4t4WZzaMJt0k9vTUJWxb3b8D9kZTflM/uMzsretrnTvqM5YN2Rbtqu8XLYEdc2melXJeGZpPuc15pvORtk6mnY8ceVWjLKh9Zv1+vQ+qzeXNql+bL2/Vc5AtX9d5/Oj/s2i1yfETztaXghdzXnLm8dQa2y+4jWdUjEdDgC4Sb6Tj4fDPgKKvUR9Di8AQWh9/wFRBHgzfB43r6g00hduMaoRrqjRy5duudPda5oq0W1qsREozRzV5A4+yK2pNDOT5q4Yc9M2VqMw5GzBaWfW47WTp2horCMlJOi1059+S9GuTO+n35G07cvjjLExPVMfhjEf6VGkaC+bj1mN/pgPm8woa4aaUlaSsWrY64hX0+w0g/SywhlTkXYvRISEqE9sXRFm3OFkefaJUsWbjwz0G7PlKea/gq65R729CgCwoiTxy2thQv7+IeHjgiUqaDaqqmkyQorVdfeaWtBVVohZ/K2SqBCz0MnUMZAQ1zUw0UEaVCOQ5jqVRnWG4qY6IrU1us/6KKnU5l9h5+6YjmgCANTwXidiMFaMltrqTsRVsgUmBOxiQjf87zFvlI+Hn/c/X9f4lTTd8O+sH9ERnn57Mf/UBeGByMAGQbj4/vPIKsW7OwuRlIhDgm5SKfxN2gIf1xPj54x2J4KjCrFYJLOUh49MwqWpSDwawsHHh4jMapGGSHJ+lUnAR7vZO6OI4Pw7qZgks95GwgpcidlmTZIASeYrCQMeImeEx1KBCBiSuS0tzHT12CtOAhPBZYNChZmlpAkzrb7pEiPMLLR2KlgszOwmSKwI4QIVZpYSIMzUra2tNxuuMuqlpPh5Yho23WEm9b1Fhm/D+JZfz7xKDN6uIQiLvKjyWvofoxIj4jqGqMTgvS5FCY8+rxKDt6PS8BCwVWJ/AYWee3a6WZK17MqF4r01ZkHR1AetC2mHtNh3lD51RsSntu2vrZHWradqm7M5f2l0q3XK2cb38qNTjyor7wWtUyqURa5pOlWoJhjmHfpUlZM/7NqtZHaRUm9r5nYT6SNPHJ848uyJP+X2WSDtNGWVQKHs7PY/Q+csEmITzj+MZrl087GR3CZrzfWZHsL8p/sQwqXO8dPMavzrTJhMT3q0rSgOsijs6b+xPs1ajCutey/L8DHmwaH7NB3MqTQ276r4MLsbCgaCDl9483FN1tpxvan+q6M13XGngtIEtBo2CPs+nD3M9lyB6pW5cYpaX+cX86/BaTbiNGwf5PsKw59zKCm7FvY7HVN+1jlZsmpP7371iNUW3UWJXq4vdrYe9up0ykYwOGTK2yhHHLmtUmqtpKkigUq+Qrf9dN9Hq4LoL4YZA8GMQ62Th5r8an0iKYQGghld3sW4z92Lz5wp2u3Q7SQQwOKTb3pHTWFSNdI7PqsgX96gdW2rjunJ0Y1Prb5+Dph/HLCrntE6QAHAACX8zwEcW/xTQgbvOtoKG2SBhKwUvoSMlLvlh8ADHktZeAgkv3rFYeD9FHjA4yQHF+NfZfVT4PHrWC0h8CDbERAfoqMH2Y6A+GFI5iYLGxRL4AGPnzKRUCRzhPOYxBJ4/LrH5JICj1K4Ag8ijnBLCjx+Hc8lBR6lyws8ugkSeEAwsQUepd8FHkbhTyTZD4/4dQlrxGrnxirFrmPcpZVNf3TgfCo6Qi4k7bap0MP7eYKecyN91e1K+XfEXDZQHKxAbxRISlszF3ezetZfZGA0xFwizadsu+uhyuOv9uX287XxfuDZWXKXO6m7rikYE61LYeO9/h5zV0jsnn5Gg90j+a4a1eqUvg6yL/84izn+JNTCvH273/rrluZ+2b4Rf/jLrhS2P2I0v61Gz/NQSYsUHkScQyEKD3JUFJeMjq3wgHfNbIMNsuR3RlWlkZGReG6vqXGFrrherUm2qTGyptKk19RYSEK3EiGeYyokLGJAM3+2yHfyLhL7fvmQKwNYOg14ad4KG2SBTgMenDIJcD91GvA22YBUzCV1GqX/STqNboJ0GnBShqNzWPo/QKfRTZBOAw7ZRTqN0v+BOo1ugnQaBNKG6DTgPQJ0yYEL0WnAexxYkhN/KZ1G6X91Gktn3o6smcej0yj9n6fT6CZIp0FsLvA9bf8X6DRIOvn81GmQ/at2GZ1G6f9unQZJR8GldBqlv1qnQVibyMffgfxtIuygZG0TPUVVEtMmav5a9ZDYNhGEC/w2EVK/zrBWz2R+hsfPmunjcpdRue91KXy7tmn5Rf2Y4AHvZCFIUGAEysHFg2yNqCWjf/8D+6Mf+3dLV0qsTV/zmUJi86lVLuuyn2cqCdbLvajRF4xqMJ4cMmWfpnjoWeTJdcKB8dPcNkdG+i56/97NuUgGf/EDnxK83yU/nQoVnFtdxMVevfbDDK9mm10HxuTiFzmqmdOrae5eecB2siPoWhaLishbWb1bASKCBiedKmsDn0cjqj+4zWRnSBzXlCxxGdox6WnHUXJVZufhNaOhppsSjyZlfewsftufgJy2E93ahClpevVhperYrYlOmmvcXopv247fcFZIjd2+5lCUPcWpkV0bzO417Jy/WSMyuULaAQApFPh2RwRG/v4Z30GOvhIugEV7ZB2ZY1gqyUozMCFqsftIaN3RRiaaa0lJRTt8rRJqyyjLVYLWV2pMK9LLCn64EScfN027suSQh/5JBY6wdTJJ1DYvpypvuPBX7kXlTB2zode+e/Rty+rBOPvkQ28ZdojXih6MUze8fYtaEtX4/CqqWGB3Op+Qrt/Ex1Zqvhjl/b+fvLXv8np6UQ03KdP3zF2mx5AzzHu6//SSLjFrj5ZXjnwiIDnIzGC7S0iJr6azR2e967aD9he2J9SM+9k/FJm+va1f5l0986zHt7CQ+RMZv07NkyoAABcF/LYtJGVkbdsuivx/vm0LyQjZ27ZLRidr2xYnwsK27TLDPS6W7zIZYOELZKOb6jj9Jux1rqCUxx+d0/tX1qrmUZdFmW67sonLJdlY4J/hHn5vCqcHDcbaf1dEK75v2q4sdaDRwmYLtUg7j55uq7jUvn3TlrcfBLXNTExPzPTUzzZ96rklxZ9o9DaY94qTWLKNUCmF88NrXgnOew3erv57uMcLRV8Kd/vaj0PXx1eXzQ/3YPg+3CPv+3CPZ7Lfh3vQ/xjuYbci9pLdWrM/ehoNQoofrt/8Wp9R/sun+m15DawyrD6vOS3oztkX0rycyygy2J0yK1rOvwV50726YIBRPr9gljVGw1o8htFkfKbuTJ++UUUzpm5nz7rjXgIqxoZSjq5d0/d87DEfxFp6fwz3uJgxcvDqNjF7qu/DPdwvaSQoCtflXpyb/+7YHC+Xa0MBgCwRnVkcu0jeziw+kOU7s72UTfRrrRMiaWj+EKa/S7/yKUp7/ty00voyVlOiTzHk9y3fL1h8iZCGt0bs4R/w7kx5ooCIPaWoEo+2+NSy5OQPqcj4x1JGcncWTv7IZx44K164O++9yFV5a9rIhJ2BIszv2I76Zn8Du39O/lh/bF1ny0D07CYDhDCLFfsukSbBH/tmNTFkuQ8AwEQBX3VHOEHssR//kESWqLGWjfi8E1YJ1wnK1uP589hmE0o37/SiqsI1L+1y72AEkSsDA55dYEnqmikfD2AeD1H2Zi7yP3qBWhfTeETzxLkpu1e074cPPL2gIsBD8+3R0XRmRXamD29P5qZFKjYUh5o4ztUJHpm4s+qLQ7PjiwbXYJkNl3sDbR81s5zw9uK+fv/ghOIFKWZM79avn9ZlovUcW4Uqtlg22p9tdiq0vPygluHZS3Yl12+YO7Vf3zntvDh7y+14nIUuBklRUDz4oEHDdv491Mm9e7QcAHCCAr5SE0YC/aVsfewDsXOoVaOu99fRqdvn3U5G9c1Wrpq3Ox074m92X9u50r557+XkJjRaeGbIVPzA6uPj+w+2fUOMq+vMylDIuhmKpV01VtG+qsuudjXr0mjTTHXJ+8GRTwd3ZVa3fb1e/Jurj6dGy5sy98Pr99Gs6A4+yZNfXiKl3cNOObLPbbjXaRD1JDErsrhYK5WZxcHG5tW0fqJ6bfMcV/Cp0gc2e48xsOtLe5YOFD/+zKdd1JHsyrLvVM7K7Cc/1SyXhVGV+QCAPrK9W8kp8sCJ8IvkEEvjkVfkgQfjX2VFXpEHHoxfJPLAh0hGkQd+mF8k8sAHSmaRx/JQv0jksTQwmUUe+ED+HZHH0isgt8gDL8ovEnlAMJcSeRyTZD88kj/EzluzyUqfp4bpuAXy6l3j/l2Nq5uMR2mDwzd2PuxfeW9uYgTl1tluxZkXBVReN7OxGTdHzIk9H5wojP7ypVhkZ7NCn2ScbS5vwGWxSZ27UfdKarOZaIUKR8fWXqlKPixRxGlzumlDKkXi1UDEIwsDl606h/0jnLrYzosbxzakIafSPjmwx2cYXqp71nTEiZf/1vD89M/hkzF8V4jUeCzgTmaNx5LRya3xwAeyZFFCNsKE/okkq9acT7g9S9KE8zj7HntezUwxgxcGA1/4r41uKXHYI214nPNPVyYnESetoQLUOcq8uo7Jr9kjOyqeT/oHol4KVhVkivZFHTJDjshVR28z2sa+91NXddrg6vDfaFXnL8q8vke3UQCAILKljsziEHwgv0AcQhgcecUhMDD/Kw6Bn7L/9eIQOGT/g8QhBNImuzgEJi7ZxSFE4v9XHEKyOITYzP8nikOIzcX/CXEIjuSQVxyCD+T/pDgEb0L+P4tDxGCseH6IB7yrZAtMiH+jFjQP6Ont4eux18+JbL/Zjh8GIfEXMXG0k7g7ygHl42PvHUgsS1EY8D5emF9Tb8VGQPj5/HUhY0FFNOQgS9S4y0ceGo10XT587OjRjeuz3qc/fr1pRNr3w/rKdPUDK2NO10k01mrfaniq2/GRVbx/dPOzgrNrsxuFmzwmD5648sZBS8q4OUAyW4kt/MgzCjuv0o2X2l+0y/LGWB0U0KYuva5GWU7BbuzhyFZdLiHALlY/0+EePG5qwPk16X4OJlmpulELvSLGOuyd9ZqXim9FclMfhd/47J9zfuep8PbKz19SDz1wT44oYagsv5Z4i6fzw7tAPk3+Qk6+rDVa1gzro2gYh1XoE/x2Gzq9ebd9S/DcH5xnapvSUjKdtP1U55yjgl9qz7JGDrNGvBl0KqfWZ8YcDbobXZk0tzOX+VzIUTSi9cb9oMQr83fkNqnPvAwUANynILGYQPCcDzpi7o35QQWkloTYF2E5e0i4efys3FXVGgnX6pobp3/9ErDyi9eZcZ8xv8w3vcj6ncbp/uNRLD5M9F/Hxm++6RVBiovViiHNderFKoUMGvXq6p8b1RsZNRkZNSFkKIX9NYVWzMurnJz3zJoAAMzx3j2rFq2NSAkHz5KBvv8h7oBBS7j9DIo5a3uTTYr16RX5K73Z46Zi905NvtAwp79L46U9ot6+yUum8rjThdy4u/Qzw7kzfNfG3bzfnh1V5OPJDo/zo/sw0Oki4ElZEJoUJ+0wa09jKKzDqyo+ENUl4HKUy2CYzeRBbWmaS/gZIUP32H1bT18KvCKIVKgf7DmBMb1h7W2hcPMAfbIMym/GJs43nBd9PWFmdkKsgJ3SmvIpH6ffDGq6eOLT6EGFMhdZRzm3GMY/5rg1bofvqHtiO5U06Wc2xVyiQTdyOFCvlnnbZdW7guuqXlpMb5lremG1Y3tNEuN6c6XQw613nokm1A9zDmypKpfnz9ARsr30wv1GiTmfHKUQvWt7qPe5zoFICaDCJ/BI5lzgkZp7zh+nmlnXbFg/dqSpI6nthlCKGGKz9o0AWWPbPS1GIbQ/9lPwiEaxNAUAz/EejNfi3gbsLQg194haW8B2KHiWfZRv2hezbggZaO61iV6Z7tQB5YaZdPaL9IGtd7p7821e6mkJdz3M5YpUVWgM7slXjroSqMQT1fR42GHu2ye03NrxoaSXxx5wnmfbtIN3X+TZVKvLTgwWYg6aDMKNrroV6/bLJQSGnpt1sQr7k522Z/2rzuAVt9Wlo+qOjYrF0InmGbpcPO/7GX1wLoUhnzf/vYYdW1WgAH/RpCTiqpyald2cQsVhayGnVmPFq2+NXczDA93iuHapo/1oTHbKXMmy4mdhvuxW+pVD3otH6BBTlkDNNmonrVcl3Bzt2Tkaoh750S4y5TG0b5/s1tUYpR/iNncQ9LDLE1yLFgrTok25K9kfdefx1nUiiOlJtao4gYT3DNGPJzflpWVqv3lsOSPIyfM2aJLtS9qHqkE3WlWvy6Vfj95e4WUyx9d+UShytY2Muv8cw62g5h6XIolLLuo7MIdKOU9svuusP8icFRvaZhZC3XDMv1zcGBnLwMJtY68dqSbM1XL7iKlAcW4+FVI5/sSGap/TNzuN1NBPD+YcHnjO0jrMsvW156D05ZS2+x5luhpa6cOI/gfKI++URv3strRzCsUZD38KLvO6mL6ranOgNqafH70jhnWHdfTtw1orB2vNQvKCGRCJHY0OjZWjfTeoVCIo3btWVrB+sTn4Mv5kF+UX+vbExxNfY3pKPj6aG5B8f1f+SFCQT8+OvRo57QPhgSaSwq12rgV0BUG0d1zG2tpSKDSYO0YPpr1YFTCmvjPqctkMONM8wFfes41255SvbYsC5x7KjIvZu8zTpp99OXVTd3TKsNc6XLqcA+l2PvTsTRZmZeMHR4SLMi5ZWyataE88MTKBlpuYPuvII3pBfN87P/7ZsnjhIpGJlUl7kdH9B1v6FexkHmqItrWjBmsx++deyokdyvmU0Tg1tPW0dpZVsxbV4TV8rWrVc9YcX/mQE8qOjz295WcbHSmmb9Nmck57UR06InJ4T4XBA2HZ1waHVQwnC0ocuQZY5z5GnFrr4pYq7pwde/dePFq0go5FXSajQaZOad/v598I2TS9Y3Jqb+/KncidngiJ0d/9cjj/2qSpbNGqD3b7WSOe73nn981thqEl2EbJ8VTlqUoEu9CFK80PZD4Yfzs6lqAQOpGc2LD31lvHrqL+a1Jz637cqZWB+YlWVAD40OJ78nIsulP9fNEY4tpVPLhiwT5D9iyIzLQw8vdAZUuujoKSnern/DpsU05usHL+780d+etPgnw/F4bEdsZcDQnpsWTIBd6eC8Nh+19yQsIpU4DlzDUXBsP2ouSCBEtdFGyxC+LPaEv7V/79DyuYU3emBHjdLH/GWcq/8mec20vFweVmCaWK7QYpCKFKTQVgu1kujI5t7CgEid6wfHT8mV3KHfJnRhyoATyvSNyBhSCBHxAUeJGBEzQt2LaHUpC0rKMBJHhF4qMhCqFhvQIQaakIZYLtJCgLYfIIHgI8MnwQMk20AKZRIZQEtooXASEhTgeINirEt35ByPpVVwIizP8I3epEwqMvMP+DZgnbQk4akqU5eAgLzf/wUZGAUHnKAIj3p4OywR7ysxXCRoQRkO5Ph4+TNIRTFjy4xVZ0uIG2QoC4mAA5rOigacS2bFOHpPE4CXD4HgZ/42L7nGlDcF2ZAZnt4RZiYxuAISHYmSzgV/iqEX6vs/4GiHQdg1LEPk9CH+ve8BAWDy/GR0YcQqYXBtRCwzHcdOQgdH5fBUg1HCOcUCEMsIVeY4QS4mYDpHqN4SMkCSEUCRvMkbBrTgHC6QsRMPBoyUJoubIDUh3GcFNThlB7TiQUPHrQt6ghByDeXAxKDFsgCH2L3oENsoS5GD5OCAinLZyABCcx3LsFJRUPHwXfRv0Nh22qJQCBY+ECcE3ECD90MqwGRLh0EXqNoQmPvtilC5ofbGurbZD81MEGgX1Qh7JKWQOI988i9JxGvxaQwz8LN5wyBC6QBLiF/lnQjcPWlBlANm6QVEyc/lkL1wCRekHWsJEb/DpDqoXLwJZeQY+MTOvBr7CKwrf/6pD9p+YBZDVtwo2rC8E9SA7cJUyboInHViBZQhI/TU58XKZNC5eDLQKygyzHfyP495yMFi4LW35jDVnWUV7w75gKEf4QvrIJEG/lA+WNLYCBHjmp+QAZrHwW4mHLV6AnC47NgHSnHHw5FIPksBwG3AKbG9xnpS3QMgo/INHmBl8dHXqWiV8eawl5Ae7o0AvuG+HRFzvc4AaRgYDsEQAkOtz83JjuBRujBNmYCuKQFnef8JHbAiGnJghINrfBTU8FQu8WsVjwCMpDCAoIATL42uBGU4SgRROPtsDXBndStSBJXSEMyOtrg48s9Hx6nATohZY2uOn+DqE7SSomPMLboAcjEUAeNxvcbDWhp3GSAOFRhdaPMKKAFCMb3AShL/MmImDg0YK2Y3XEAOEeNlAS2B4w0HZs4rJBFw8nWxgd29EF2o5FiwP4rjL4MgLt0H1bPvwCqRp04bg7dGYSgGiTGMIPGCgEIMJ4BXd0aNW+hPDoCyb4QLOE7dAB7dDxSwJSjFcIPyvRSQPi7T+gbHB36NxggywxZIbw7RmBAQcdAoH7AQVtqpjJANKcPvCRgfYgSmBALZz/gJsOtAchKAtINfkgnFAcDLCFox8IJTQDF4OUl6LHFkCKtQehL8UWImDg0YI2VfTkAKmuHripQZsqd4mEgkcP+uBDyAPiDT2gxHA3VaJhgywx6wEfJ2hThXMrIMG9A/duQUkdhI+Cb6P+hsP2sYA2VYaXh1s404Hwc8KgAiDCGIPQa0xvGyDeGAOaH2wfCmhTJRs2COyzFZSVgSIg3uwCNwi0qZILG2TJeQaEfrTyKgGymV1ANw53U+UiqZg45xksXAPupkqFMvh17hELl4G7qZKvAn6FrwO+/Yc2VW6pArI6LODGhTZVNquBX+KwAE087qZKMjnxcf0S/cLl4G6q8GiAf892YOGycDdVRDTBv+MAQPhDeLsWIH7uPpQ37qbKLdggS/1q9UI83E2Vh9qA9LH2hJd9PH4HhM+kh3LAnukOLfuULht0+bIP9nB2aNnn7HYAf0o8voxAVfj8OgDevHdCz1mHCAq8cJIqNCvYs8uhNaVXhEdfNDkVHwcRCAduJCBuKjqUBva4b2jRJxgWACwm0Opeth6AOUuc0OoevT4gepY44VcTtwEgYj43oeWjA4RHx1vdw57yDN3oV/AQiK/uXTMCxI+QhrLBHtUMfWrTGgPSR0jj4wTV39+EB7d4GjShPUJ2E0CGadDQLGIPTlaFCjeIR8PzJPgbFXvaMLSV5GAKyDqieSEy9tReaMsu2QyQe7Yx4bd4yw5A5JhgQivEW80BaWOCCS+opsGAIrZCzGgBSJ0QTDihQzDAiK0Q98LFIKVCvGMnIGUuMKEV4vtEwJBSIZa0BKSOBCa0QhxPJBQpFWJ2K0D8NGBCK8SBsEFIqhBPwYAjvkJsvwuQPvp3US8cawoutEL8fHk44ivEz3YDIsbqEtyFsAbEj9WF5gd7GC20QhwNG4TECrGrDSB+4i2hFeJ62CAkVYhV9gCyTbwltEKcSyomGSrEb2zBrxshS3iFuNEO/IrhroRXiCvtAVnHrBJaIVbfC37JmFVCK8SPyYlPeoVY0RH8e7NHCa8Qb0eBf2cMKOEPYUcnQPzwTUIrxJWwQUirELc7A9JnW+LLIVR2H+ECiB1MifusBJXdD8GEIOU7xA4NSJlJSag+tZwIGHwzKQn/rv/dFRA5jhLKDXtIIvS7Pg0ewuJxlISfPHkwgIgZjLg3SRKq4iA8+vJnc+xxh+wQmP5FMAvGL+LLxypIPkLdAP5pidA1YY/s44GsqXbJQEtOS1wYFHu62FpI0B0egJDZf/jockDoznqDZUeU4d5sKGOkDyB0RBn2ApeaR/ZzgckLg2JNJ6NZ8dffKgAFIFIfgH6fv/7r/wUAAP//sbsHCIzAAAA="); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}

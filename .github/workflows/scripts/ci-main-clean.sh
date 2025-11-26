#!/usr/bin/env bash

dirpath=$1

# Extract the base directory name for pattern matching
if [ -n "$dirpath" ]; then
    dirname=$(basename "$dirpath")
    echo "Cleaning Docker resources for path: $dirpath (pattern: $dirname)"
    df -h /
    
    # Process containers and images based on the directory
    case "$dirname" in
        # "mysql")
        #     echo "Cleaning mysql resources..."
        #     containers=$(docker ps -aq --filter "name=$dirname" 2>/dev/null)
        #     if [ -n "$containers" ]; then
        #         echo "Stopping and removing mysql containers..."
        #         docker stop $containers 2>/dev/null || true
        #         docker rm -f $containers 2>/dev/null || true
        #     fi
        #     docker rmi -f $(docker images -q mysql 2>/dev/null) 2>/dev/null || true
        #     ;;
        "mssql")
            echo "Cleaning mssql resources..."
            containers=$(docker ps -aq --filter "name=$dirname" 2>/dev/null)
            if [ -n "$containers" ]; then
                echo "Stopping and removing mssql containers..."
                docker stop $containers 2>/dev/null || true
                docker rm -f $containers 2>/dev/null || true
            fi
            docker rmi -f $(docker images -q mcr.microsoft.com/mssql/server 2>/dev/null) 2>/dev/null || true
        ;;
        "pgsql")
            echo "Cleaning postgres resources..."
            containers=$(docker ps -aq --filter "name=$dirname" 2>/dev/null)
            if [ -n "$containers" ]; then
                echo "Stopping and removing postgres containers..."
                docker stop $containers 2>/dev/null || true
                docker rm -f $containers 2>/dev/null || true
            fi
            docker rmi -f $(docker images -q postgres 2>/dev/null) 2>/dev/null || true
        ;;
        "opengauss")
            echo "Cleaning opengauss resources..."
            containers=$(docker ps -aq --filter "name=$dirname" 2>/dev/null)
            if [ -n "$containers" ]; then
                echo "Stopping and removing opengauss containers..."
                docker stop $containers 2>/dev/null || true
                docker rm -f $containers 2>/dev/null || true
            fi
            docker rmi -f $(docker images -q opengauss 2>/dev/null) 2>/dev/null || true
        ;;
        "oracle")
            echo "Cleaning oracle resources..."
            containers=$(docker ps -aq --filter "name=$dirname" 2>/dev/null)
            if [ -n "$containers" ]; then
                echo "Stopping and removing oracle containers..."
                docker stop $containers 2>/dev/null || true
                docker rm -f $containers 2>/dev/null || true
            fi
            docker rmi -f $(docker images -q loads/oracle-xe-11g-r2 2>/dev/null) 2>/dev/null || true
        ;;
        "dm")
            echo "Cleaning dm resources..."
            containers=$(docker ps -aq --filter "name=$dirname" 2>/dev/null)
            if [ -n "$containers" ]; then
                echo "Stopping and removing dm containers..."
                docker stop $containers 2>/dev/null || true
                docker rm -f $containers 2>/dev/null || true
            fi
            docker rmi -f $(docker images -q loads/dm 2>/dev/null) 2>/dev/null || true
        ;;
        "clickhouse")
            echo "Cleaning clickhouse resources..."
            containers=$(docker ps -aq --filter "name=$dirname" 2>/dev/null)
            if [ -n "$containers" ]; then
                echo "Stopping and removing clickhouse containers..."
                docker stop $containers 2>/dev/null || true
                docker rm -f $containers 2>/dev/null || true
            fi
            docker rmi -f $(docker images -q clickhouse/clickhouse-server 2>/dev/null) 2>/dev/null || true
        ;;
        # "redis")
        #     echo "Cleaning redis resources..."
        #     containers=$(docker ps -aq --filter "name=$dirname" 2>/dev/null)
        #     if [ -n "$containers" ]; then
        #         echo "Stopping and removing redis containers..."
        #         docker stop $containers 2>/dev/null || true
        #         docker rm -f $containers 2>/dev/null || true
        #     fi
        #     docker rmi -f $(docker images -q redis loads/redis loads/redis-sentinel 2>/dev/null) 2>/dev/null || true
        #     ;;
        "etcd")
            echo "Cleaning etcd resources..."
            containers=$(docker ps -aq --filter "name=$dirname" 2>/dev/null)
            if [ -n "$containers" ]; then
                echo "Stopping and removing etcd containers..."
                docker stop $containers 2>/dev/null || true
                docker rm -f $containers 2>/dev/null || true
            fi
            docker rmi -f $(docker images -q bitnamilegacy/etcd 2>/dev/null) 2>/dev/null || true
        ;;
        # "consul")
        #     echo "Cleaning consul resources..."
        #     containers=$(docker ps -aq --filter "name=$dirname" 2>/dev/null)
        #     if [ -n "$containers" ]; then
        #         echo "Stopping and removing consul containers..."
        #         docker stop $containers 2>/dev/null || true
        #         docker rm -f $containers 2>/dev/null || true
        #     fi
        #     docker rmi -f $(docker images -q consul 2>/dev/null) 2>/dev/null || true
        #     ;;
        # "nacos")
        #     echo "Cleaning nacos resources..."
        #     containers=$(docker ps -aq --filter "name=$dirname" 2>/dev/null)
        #     if [ -n "$containers" ]; then
        #         echo "Stopping and removing nacos containers..."
        #         docker stop $containers 2>/dev/null || true
        #         docker rm -f $containers 2>/dev/null || true
        #     fi
        #     docker rmi -f $(docker images -q nacos/nacos-server 2>/dev/null) 2>/dev/null || true
        #     ;;
        # "polaris")
        #     echo "Cleaning polaris resources..."
        #     containers=$(docker ps -aq --filter "name=$dirname" 2>/dev/null)
        #     if [ -n "$containers" ]; then
        #         echo "Stopping and removing polaris containers..."
        #         docker stop $containers 2>/dev/null || true
        #         docker rm -f $containers 2>/dev/null || true
        #     fi
        #     docker rmi -f $(docker images -q polarismesh/polaris-standalone 2>/dev/null) 2>/dev/null || true
        #     ;;
        "zookeeper")
            echo "Cleaning zookeeper resources..."
            containers=$(docker ps -aq --filter "name=$dirname" 2>/dev/null)
            if [ -n "$containers" ]; then
                echo "Stopping and removing zookeeper containers..."
                docker stop $containers 2>/dev/null || true
                docker rm -f $containers 2>/dev/null || true
            fi
            docker rmi -f $(docker images -q zookeeper 2>/dev/null) 2>/dev/null || true
        ;;
        # "apollo")
        #     echo "Cleaning apollo resources..."
        #     containers=$(docker ps -aq --filter "name=$dirname" 2>/dev/null)
        #     if [ -n "$containers" ]; then
        #         echo "Stopping and removing apollo containers..."
        #         docker stop $containers 2>/dev/null || true
        #         docker rm -f $containers 2>/dev/null || true
        #     fi
        #     docker rmi -f $(docker images -q loads/apollo-quick-start 2>/dev/null) 2>/dev/null || true
        #     ;;
        *)
            # No matching pattern, skip cleanup
            echo "No specific Docker cleanup rule for '$dirname', skipping cleanup"
        ;;
    esac
    
    # Remove dangling images and volumes to free up space
    echo "Removing dangling images and unused volumes..."
    docker image prune -f 2>/dev/null || true
    docker volume prune -f 2>/dev/null || true
    
    echo "Docker cleanup completed for $dirname"
    docker system df
    df -h /
fi

# df -h /
# Filesystem      Size  Used Avail Use% Mounted on
# /dev/root        72G   67G  5.4G  93% /
# tmpfs           7.9G   84K  7.9G   1% /dev/shm
# tmpfs           3.2G  2.6M  3.2G   1% /run
# tmpfs           5.0M     0  5.0M   0% /run/lock
# /dev/sdb16      881M   62M  758M   8% /boot
# /dev/sdb15      105M  6.2M   99M   6% /boot/efi
# /dev/sda1        74G  4.1G   66G   6% /mnt
# tmpfs           1.6G   12K  1.6G   1% /run/user/1001

# runner@runnervmg1sw1:~/work/gf/gf$ docker system df
# TYPE            TOTAL     ACTIVE    SIZE      RECLAIMABLE
# Images          18        11        8.326GB   1.644GB (19%)
# Containers      11        11        2.692GB   0B (0%)
# Local Volumes   11        8         665.7MB   211.9MB (31%)
# Build Cache     0         0         0B        0B

# runner@runnervmg1sw1:~/work/gf/gf$ docker images
# REPOSITORY                       TAG                               IMAGE ID       CREATED         SIZE
# alpine/curl                      latest                            99fd43792a61   2 days ago      13.5MB
# postgres                         17-alpine                         b6bf692a8125   9 days ago      278MB
# zookeeper                        3.8                               2f26c02b94ca   10 days ago     306MB
# mariadb                          11.4                              063fb6684f96   10 days ago     332MB
# mcr.microsoft.com/mssql/server   2022-latest                       a2fbff321505   4 weeks ago     1.61GB
# clickhouse/clickhouse-server     24.11.1.2557-alpine               2eee9fd3ae74   12 months ago   539MB
# redis                            7.0                               7705dd2858c1   18 months ago   109MB
# consul                           1.15                              686495461132   20 months ago   155MB
# mysql                            5.7                               5107333e08a8   23 months ago   501MB
# polarismesh/polaris-standalone   v1.17.2                           b7a8cf0a8438   2 years ago     545MB
# bitnamilegacy/etcd               3.4.24                            74ae5e205ac5   2 years ago     134MB
# nacos/nacos-server               v2.1.2                            a978644d9246   2 years ago     1.06GB
# loads/redis                      7.0-sentinel                      6f12d40540ba   3 years ago     114MB
# loads/dm                         v8.1.2.128_ent_x86_64_ctm_pack4   ccb727ce9dce   3 years ago     432MB
# loads/redis-sentinel             7.0                               6818c626f5ca   3 years ago     104MB
# loads/apollo-quick-start         latest                            8490de672148   3 years ago     190MB
# alpine                           3.8                               c8bccc0af957   5 years ago     4.41MB
# loads/oracle-xe-11g-r2           11.2.0                            0d19fd2e072e   6 years ago     2.1GB

# runner@runnervmg1sw1:~/work/gf/gf$ docker ps -s
# CONTAINER ID   IMAGE                                              COMMAND                  CREATED         STATUS                   PORTS                                                                                                                                                                                                                         NAMES                                                                               SIZE
# 8214f83420c6   zookeeper:3.8                                      "/docker-entrypoint.…"   6 minutes ago   Up 6 minutes             2888/tcp, 3888/tcp, 0.0.0.0:2181->2181/tcp, [::]:2181->2181/tcp, 8080/tcp                                                                                                                                                     d66bac92ae9646f688f70ed4b5176f14_zookeeper38_3a22ef                                 33kB (virtual 306MB)
# 8938d73842e8   loads/dm:v8.1.2.128_ent_x86_64_ctm_pack4           "/bin/bash /opt/star…"   6 minutes ago   Up 6 minutes             0.0.0.0:5236->5236/tcp, [::]:5236->5236/tcp                                                                                                                                                                                   ca280fbdb86f40c2acf86d7d526c6285_loadsdmv812128_ent_x86_64_ctm_pack4_770a59         844MB (virtual 1.28GB)
# 0d3a653fe1f2   loads/oracle-xe-11g-r2:11.2.0                      "/bin/sh -c '/usr/sb…"   6 minutes ago   Up 6 minutes             22/tcp, 8080/tcp, 0.0.0.0:1521->1521/tcp, [::]:1521->1521/tcp                                                                                                                                                                 2048856d428c4967b1c35193eb8c9192_loadsoraclexe11gr21120_295d54                      1.3GB (virtual 3.4GB)
# ca3936189166   polarismesh/polaris-standalone:v1.17.2             "/bin/bash run.sh"       6 minutes ago   Up 6 minutes             0.0.0.0:8090-8091->8090-8091/tcp, [::]:8090-8091->8090-8091/tcp, 8080/tcp, 8100-8101/tcp, 0.0.0.0:8093->8093/tcp, [::]:8093->8093/tcp, 8761/tcp, 15010/tcp, 0.0.0.0:9090-9091->9090-9091/tcp, [::]:9090-9091->9090-9091/tcp   cbd43dceef754e2d8aab507e33167be7_polarismeshpolarisstandalonev1172_ca40b6           299MB (virtual 844MB)
# 26169dad485e   clickhouse/clickhouse-server:24.11.1.2557-alpine   "/entrypoint.sh"         6 minutes ago   Up 6 minutes             0.0.0.0:8123->8123/tcp, [::]:8123->8123/tcp, 0.0.0.0:9000-9001->9000-9001/tcp, [::]:9000-9001->9000-9001/tcp, 9009/tcp                                                                                                        f1c7766fbe36401792a6f735d7acf123_clickhouseclickhouseserver241112557alpine_cfc034   338kB (virtual 539MB)
# 04689a1d581f   mcr.microsoft.com/mssql/server:2022-latest         "/opt/mssql/bin/laun…"   6 minutes ago   Up 6 minutes (healthy)   0.0.0.0:1433->1433/tcp, [::]:1433->1433/tcp                                                                                                                                                                                   41d685349a7640b28230db8d0f60efe7_mcrmicrosoftcommssqlserver2022latest_fe29fb        108MB (virtual 1.72GB)
# d5fbc5f811af   postgres:17-alpine                                 "docker-entrypoint.s…"   6 minutes ago   Up 6 minutes (healthy)   0.0.0.0:5432->5432/tcp, [::]:5432->5432/tcp                                                                                                                                                                                   2783be71b5ce417ab9a31428e7b4d8f2_postgres17alpine_c60840                            63B (virtual 278MB)
# da96a7ad7a01   mariadb:11.4                                       "docker-entrypoint.s…"   7 minutes ago   Up 7 minutes             0.0.0.0:3307->3306/tcp, [::]:3307->3306/tcp                                                                                                                                                                                   45eed646fa6c4a698893ee11cda95a4c_mariadb114_3a9cd6                                  2B (virtual 332MB)
# 27ba1904ba3a   mysql:5.7                                          "docker-entrypoint.s…"   7 minutes ago   Up 7 minutes             0.0.0.0:3306->3306/tcp, [::]:3306->3306/tcp, 33060/tcp                                                                                                                                                                        ea6d7a4c207d427a95b5ae0db91fdf56_mysql57_c21053                                     4B (virtual 501MB)
# 518e785d1bb6   redis:7.0                                          "docker-entrypoint.s…"   7 minutes ago   Up 7 minutes (healthy)   0.0.0.0:6379->6379/tcp, [::]:6379->6379/tcp                                                                                                                                                                                   af6044fc849e441bbc6c48f7a5ec5fec_redis70_b11994                                     0B (virtual 109MB)
# 7495ec2cd8e3   bitnamilegacy/etcd:3.4.24                          "/opt/bitnami/script…"   7 minutes ago   Up 7 minutes             0.0.0.0:2379->2379/tcp, [::]:2379->2379/tcp, 2380/tcp                                                                                                                                                                         49f2a2a6bf3a4fae842cc950bbc3658a_bitnamilegacyetcd3424_1265e1                       145MB (virtual 279MB)
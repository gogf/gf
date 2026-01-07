#!/bin/bash
#
# GoFrame Docker Services Manager
# 用于本地开发时管理测试用的Docker服务
#

set -e

# 容器名前缀
PREFIX="goframe"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 服务定义
declare -A SERVICES
declare -A SERVICE_PORTS
declare -A SERVICE_ENVS
declare -A SERVICE_OPTS

# 基础服务
SERVICES["etcd"]="bitnamilegacy/etcd:3.4.24"
SERVICE_PORTS["etcd"]="2379:2379"
SERVICE_ENVS["etcd"]="-e ALLOW_NONE_AUTHENTICATION=yes"

SERVICES["redis"]="redis:7.0"
SERVICE_PORTS["redis"]="6379:6379"
SERVICE_OPTS["redis"]="--health-cmd 'redis-cli ping' --health-interval 10s --health-timeout 5s --health-retries 5"

SERVICES["mysql"]="mysql:5.7"
SERVICE_PORTS["mysql"]="3306:3306"
SERVICE_ENVS["mysql"]="-e MYSQL_DATABASE=test -e MYSQL_ROOT_PASSWORD=12345678"

SERVICES["mariadb"]="mariadb:11.4"
SERVICE_PORTS["mariadb"]="3307:3306"
SERVICE_ENVS["mariadb"]="-e MARIADB_DATABASE=test -e MARIADB_ROOT_PASSWORD=12345678"

SERVICES["postgres"]="postgres:17-alpine"
SERVICE_PORTS["postgres"]="5432:5432"
SERVICE_ENVS["postgres"]="-e POSTGRES_PASSWORD=12345678 -e POSTGRES_USER=postgres -e POSTGRES_DB=test -e TZ=Asia/Shanghai"
SERVICE_OPTS["postgres"]="--health-cmd pg_isready --health-interval 10s --health-timeout 5s --health-retries 5"

SERVICES["mssql"]="mcr.microsoft.com/mssql/server:2022-latest"
SERVICE_PORTS["mssql"]="1433:1433"
SERVICE_ENVS["mssql"]="-e TZ=Asia/Shanghai -e ACCEPT_EULA=Y -e MSSQL_SA_PASSWORD=LoremIpsum86"

SERVICES["clickhouse"]="clickhouse/clickhouse-server:24.11.1.2557-alpine"
SERVICE_PORTS["clickhouse"]="9000:9000 -p 8123:8123 -p 9001:9001"

SERVICES["polaris"]="polarismesh/polaris-standalone:v1.17.2"
SERVICE_PORTS["polaris"]="8090:8090 -p 8091:8091 -p 8093:8093 -p 9090:9090 -p 9091:9091"

SERVICES["oracle"]="loads/oracle-xe-11g-r2:11.2.0"
SERVICE_PORTS["oracle"]="1521:1521"
SERVICE_ENVS["oracle"]="-e ORACLE_ALLOW_REMOTE=true -e ORACLE_SID=XE -e ORACLE_DB_USER_NAME=system -e ORACLE_DB_PASSWORD=oracle"

SERVICES["dm"]="loads/dm:v8.1.2.128_ent_x86_64_ctm_pack4"
SERVICE_PORTS["dm"]="5236:5236"

SERVICES["gaussdb"]="opengauss/opengauss:7.0.0-RC1.B023"
SERVICE_PORTS["gaussdb"]="9950:5432"
SERVICE_ENVS["gaussdb"]="-e GS_PASSWORD=UTpass@1234 -e TZ=Asia/Shanghai"
SERVICE_OPTS["gaussdb"]="--privileged=true"

SERVICES["zookeeper"]="zookeeper:3.8"
SERVICE_PORTS["zookeeper"]="2181:2181"

# 服务分组
GROUP_DB="mysql mariadb postgres mssql oracle dm gaussdb clickhouse"
GROUP_CACHE="redis etcd"
GROUP_REGISTRY="polaris zookeeper"
GROUP_ALL="etcd redis mysql mariadb postgres mssql clickhouse polaris oracle dm gaussdb zookeeper"

# 工作目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
WORKFLOW_DIR="$PROJECT_ROOT/.github/workflows"

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[OK]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Docker是否可用
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安装或不在PATH中"
        exit 1
    fi
    if ! docker info &> /dev/null; then
        print_error "Docker 服务未运行"
        exit 1
    fi
}

# 获取容器名
get_container_name() {
    echo "${PREFIX}-$1"
}

# 启动单个服务
start_service() {
    local service=$1
    local container_name=$(get_container_name "$service")
    local image="${SERVICES[$service]}"
    local ports="${SERVICE_PORTS[$service]}"
    local envs="${SERVICE_ENVS[$service]}"
    local opts="${SERVICE_OPTS[$service]}"

    if [ -z "$image" ]; then
        print_error "未知服务: $service"
        return 1
    fi

    # 检查容器是否已存在
    if docker ps -a --format '{{.Names}}' | grep -q "^${container_name}$"; then
        if docker ps --format '{{.Names}}' | grep -q "^${container_name}$"; then
            print_warning "$service 已在运行"
            return 0
        else
            print_info "启动已存在的容器 $service..."
            docker start "$container_name" > /dev/null
            print_success "$service 已启动"
            return 0
        fi
    fi

    print_info "启动 $service..."
    
    # 构建docker run命令
    local cmd="docker run -d --name $container_name"
    
    # 添加端口映射
    for port in $ports; do
        cmd="$cmd -p $port"
    done
    
    # 添加环境变量
    if [ -n "$envs" ]; then
        cmd="$cmd $envs"
    fi
    
    # 添加其他选项
    if [ -n "$opts" ]; then
        cmd="$cmd $opts"
    fi
    
    cmd="$cmd $image"
    
    if eval "$cmd" > /dev/null 2>&1; then
        print_success "$service 已启动 (容器: $container_name)"
    else
        print_error "$service 启动失败"
        return 1
    fi
}

# 停止单个服务
stop_service() {
    local service=$1
    local container_name=$(get_container_name "$service")

    if docker ps --format '{{.Names}}' | grep -q "^${container_name}$"; then
        print_info "停止 $service..."
        docker stop "$container_name" > /dev/null
        print_success "$service 已停止"
    else
        print_warning "$service 未在运行"
    fi
}

# 删除单个服务
remove_service() {
    local service=$1
    local container_name=$(get_container_name "$service")

    if docker ps -a --format '{{.Names}}' | grep -q "^${container_name}$"; then
        print_info "删除 $service..."
        docker rm -f "$container_name" > /dev/null
        print_success "$service 已删除"
    else
        print_warning "$service 容器不存在"
    fi
}

# 查看服务日志
logs_service() {
    local service=$1
    local container_name=$(get_container_name "$service")
    local lines=${2:-100}

    if docker ps -a --format '{{.Names}}' | grep -q "^${container_name}$"; then
        docker logs --tail "$lines" -f "$container_name"
    else
        print_error "$service 容器不存在"
        return 1
    fi
}

# 启动docker-compose服务
start_compose_service() {
    local service=$1
    local compose_file=""
    
    case $service in
        apollo)
            compose_file="$WORKFLOW_DIR/apollo/docker-compose.yml"
            ;;
        nacos)
            compose_file="$WORKFLOW_DIR/nacos/docker-compose.yml"
            ;;
        redis-cluster)
            compose_file="$WORKFLOW_DIR/redis/docker-compose.yml"
            ;;
        consul)
            compose_file="$WORKFLOW_DIR/consul/docker-compose.yml"
            ;;
        *)
            print_error "未知compose服务: $service"
            return 1
            ;;
    esac
    
    if [ -f "$compose_file" ]; then
        print_info "启动 $service (docker-compose)..."
        docker compose -f "$compose_file" up -d
        print_success "$service 已启动"
    else
        print_error "compose文件不存在: $compose_file"
        return 1
    fi
}

# 停止docker-compose服务
stop_compose_service() {
    local service=$1
    local compose_file=""
    
    case $service in
        apollo)
            compose_file="$WORKFLOW_DIR/apollo/docker-compose.yml"
            ;;
        nacos)
            compose_file="$WORKFLOW_DIR/nacos/docker-compose.yml"
            ;;
        redis-cluster)
            compose_file="$WORKFLOW_DIR/redis/docker-compose.yml"
            ;;
        consul)
            compose_file="$WORKFLOW_DIR/consul/docker-compose.yml"
            ;;
        *)
            print_error "未知compose服务: $service"
            return 1
            ;;
    esac
    
    if [ -f "$compose_file" ]; then
        print_info "停止 $service (docker-compose)..."
        docker compose -f "$compose_file" down
        print_success "$service 已停止"
    else
        print_error "compose文件不存在: $compose_file"
        return 1
    fi
}

# 显示服务状态
show_status() {
    echo ""
    echo -e "${CYAN}========== GoFrame Docker Services Status ==========${NC}"
    echo ""
    printf "%-15s %-12s %-30s %s\n" "SERVICE" "STATUS" "CONTAINER" "PORTS"
    echo "--------------------------------------------------------------------------------"
    
    for service in $GROUP_ALL; do
        local container_name=$(get_container_name "$service")
        local status="stopped"
        local ports="-"
        
        if docker ps --format '{{.Names}}' 2>/dev/null | grep -q "^${container_name}$"; then
            status="${GREEN}running${NC}"
            ports=$(docker port "$container_name" 2>/dev/null | tr '\n' ' ' || echo "-")
        elif docker ps -a --format '{{.Names}}' 2>/dev/null | grep -q "^${container_name}$"; then
            status="${YELLOW}stopped${NC}"
        else
            status="${RED}not created${NC}"
        fi
        
        printf "%-15s %-22b %-30s %s\n" "$service" "$status" "$container_name" "$ports"
    done
    
    echo ""
    echo -e "${CYAN}========== Compose Services ==========${NC}"
    echo ""
    
    for compose_svc in apollo nacos redis-cluster consul; do
        local running=0
        case $compose_svc in
            apollo)
                running=$(docker ps --filter "name=apollo" --format '{{.Names}}' 2>/dev/null | wc -l)
                ;;
            nacos)
                running=$(docker ps --filter "name=nacos" --format '{{.Names}}' 2>/dev/null | wc -l)
                ;;
            redis-cluster)
                running=$(docker ps --filter "name=redis-" --format '{{.Names}}' 2>/dev/null | wc -l)
                ;;
            consul)
                running=$(docker ps --filter "name=consul" --format '{{.Names}}' 2>/dev/null | wc -l)
                ;;
        esac
        
        if [ "$running" -gt 0 ]; then
            printf "%-15s ${GREEN}running${NC} (%d containers)\n" "$compose_svc" "$running"
        else
            printf "%-15s ${RED}stopped${NC}\n" "$compose_svc"
        fi
    done
    
    echo ""
}

# 显示服务信息
show_service_info() {
    echo ""
    echo -e "${CYAN}========== Available Services ==========${NC}"
    echo ""
    echo -e "${YELLOW}基础服务 (独立容器):${NC}"
    echo ""
    printf "%-15s %-50s %s\n" "SERVICE" "IMAGE" "PORTS"
    echo "--------------------------------------------------------------------------------"
    
    for service in $GROUP_ALL; do
        printf "%-15s %-50s %s\n" "$service" "${SERVICES[$service]}" "${SERVICE_PORTS[$service]}"
    done
    
    echo ""
    echo -e "${YELLOW}Compose服务 (多容器):${NC}"
    echo "  apollo        - Apollo配置中心 (8080, 8070, 8060, 13306)"
    echo "  nacos         - Nacos注册中心 (8848, 9848, 9555)"
    echo "  redis-cluster - Redis主从+哨兵集群 (6380-6382, 26379-26381)"
    echo "  consul        - Consul服务发现 (8500, 8600)"
    echo ""
    echo -e "${YELLOW}服务分组:${NC}"
    echo "  db       - 数据库: $GROUP_DB"
    echo "  cache    - 缓存: $GROUP_CACHE"
    echo "  registry - 注册中心: $GROUP_REGISTRY"
    echo "  all      - 所有基础服务"
    echo ""
}

# 显示帮助
show_help() {
    echo ""
    echo -e "${CYAN}GoFrame Docker Services Manager${NC}"
    echo ""
    echo "用法: $0 <command> [service|group] [options]"
    echo ""
    echo "命令:"
    echo "  start <service|group>    启动服务或服务组"
    echo "  stop <service|group>     停止服务或服务组"
    echo "  restart <service|group>  重启服务或服务组"
    echo "  remove <service|group>   删除服务容器"
    echo "  logs <service> [lines]   查看服务日志 (默认100行)"
    echo "  status                   显示所有服务状态"
    echo "  info                     显示可用服务信息"
    echo "  clean                    删除所有goframe容器"
    echo "  pull [service]           拉取镜像"
    echo ""
    echo "服务:"
    echo "  基础服务: etcd, redis, mysql, mariadb, postgres, mssql,"
    echo "           clickhouse, polaris, oracle, dm, gaussdb, zookeeper"
    echo "  Compose: apollo, nacos, redis-cluster, consul"
    echo ""
    echo "服务组:"
    echo "  db       - 所有数据库服务"
    echo "  cache    - 缓存服务 (redis, etcd)"
    echo "  registry - 注册中心 (polaris, zookeeper)"
    echo "  all      - 所有基础服务"
    echo ""
    echo "示例:"
    echo "  $0 start mysql           # 启动MySQL"
    echo "  $0 start db              # 启动所有数据库"
    echo "  $0 start all             # 启动所有基础服务"
    echo "  $0 start apollo          # 启动Apollo (compose)"
    echo "  $0 stop all              # 停止所有基础服务"
    echo "  $0 logs mysql 50         # 查看MySQL最近50行日志"
    echo "  $0 status                # 查看服务状态"
    echo ""
}

# 解析服务组
parse_services() {
    local input=$1
    case $input in
        db)
            echo "$GROUP_DB"
            ;;
        cache)
            echo "$GROUP_CACHE"
            ;;
        registry)
            echo "$GROUP_REGISTRY"
            ;;
        all)
            echo "$GROUP_ALL"
            ;;
        *)
            echo "$input"
            ;;
    esac
}

# 判断是否为compose服务
is_compose_service() {
    local service=$1
    case $service in
        apollo|nacos|redis-cluster|consul)
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

# 拉取镜像
pull_images() {
    local services=$1
    
    if [ -z "$services" ]; then
        services="$GROUP_ALL"
    fi
    
    for service in $services; do
        if [ -n "${SERVICES[$service]}" ]; then
            print_info "拉取镜像: ${SERVICES[$service]}"
            docker pull "${SERVICES[$service]}"
        fi
    done
}

# 清理所有goframe容器
clean_all() {
    print_info "删除所有 $PREFIX 容器..."
    local containers=$(docker ps -a --filter "name=$PREFIX" --format '{{.Names}}')
    
    if [ -n "$containers" ]; then
        for container in $containers; do
            docker rm -f "$container" > /dev/null
            print_success "已删除: $container"
        done
    else
        print_info "没有找到 $PREFIX 容器"
    fi
}

# 获取服务状态标记
get_service_status_mark() {
    local service=$1
    local container_name=$(get_container_name "$service")
    
    if docker ps --format '{{.Names}}' 2>/dev/null | grep -q "^${container_name}$"; then
        echo -e "${GREEN}*${NC}"
    else
        echo " "
    fi
}

# 获取compose服务状态标记
get_compose_status_mark() {
    local service=$1
    local running=0
    
    case $service in
        apollo)
            running=$(docker ps --filter "name=apollo" --format '{{.Names}}' 2>/dev/null | wc -l)
            ;;
        nacos)
            running=$(docker ps --filter "name=nacos" --format '{{.Names}}' 2>/dev/null | wc -l)
            ;;
        redis-cluster)
            running=$(docker ps --filter "name=redis-" --format '{{.Names}}' 2>/dev/null | wc -l)
            ;;
        consul)
            running=$(docker ps --filter "name=consul" --format '{{.Names}}' 2>/dev/null | wc -l)
            ;;
    esac
    
    if [ "$running" -gt 0 ]; then
        echo -e "${GREEN}*${NC}"
    else
        echo " "
    fi
}

# 服务选择菜单
select_service_menu() {
    local action=$1
    local action_name=$2
    
    echo ""
    echo -e "${CYAN}========== 选择${action_name}的服务 ==========${NC}"
    
    # 停止/重启/日志操作时显示运行状态
    if [[ "$action" == "stop" || "$action" == "restart" || "$action" == "logs" ]]; then
        echo -e "  (${GREEN}*${NC} 表示正在运行)"
    fi
    echo ""
    echo -e "${YELLOW}基础服务:${NC}"
    printf "  %b1) etcd         %b2) redis        %b3) mysql\n" \
        "$(get_service_status_mark etcd)" "$(get_service_status_mark redis)" "$(get_service_status_mark mysql)"
    printf "  %b4) mariadb      %b5) postgres     %b6) mssql\n" \
        "$(get_service_status_mark mariadb)" "$(get_service_status_mark postgres)" "$(get_service_status_mark mssql)"
    printf "  %b7) clickhouse   %b8) polaris      %b9) oracle\n" \
        "$(get_service_status_mark clickhouse)" "$(get_service_status_mark polaris)" "$(get_service_status_mark oracle)"
    printf " %b10) dm          %b11) gaussdb     %b12) zookeeper\n" \
        "$(get_service_status_mark dm)" "$(get_service_status_mark gaussdb)" "$(get_service_status_mark zookeeper)"
    echo ""
    echo -e "${YELLOW}Compose服务:${NC}"
    printf " %b13) apollo      %b14) nacos       %b15) redis-cluster\n" \
        "$(get_compose_status_mark apollo)" "$(get_compose_status_mark nacos)" "$(get_compose_status_mark redis-cluster)"
    printf " %b16) consul\n" "$(get_compose_status_mark consul)"
    echo ""
    echo -e "${YELLOW}服务组:${NC}"
    echo "  17) db (所有数据库)    18) cache (缓存服务)"
    echo "  19) registry (注册中心) 20) all (所有基础服务)"
    echo ""
    echo "   0) 返回上级菜单"
    echo ""
    read -p "请选择 [0-20]: " svc_choice
    
    local svc=""
    case $svc_choice in
        1) svc="etcd" ;;
        2) svc="redis" ;;
        3) svc="mysql" ;;
        4) svc="mariadb" ;;
        5) svc="postgres" ;;
        6) svc="mssql" ;;
        7) svc="clickhouse" ;;
        8) svc="polaris" ;;
        9) svc="oracle" ;;
        10) svc="dm" ;;
        11) svc="gaussdb" ;;
        12) svc="zookeeper" ;;
        13) svc="apollo" ;;
        14) svc="nacos" ;;
        15) svc="redis-cluster" ;;
        16) svc="consul" ;;
        17) svc="db" ;;
        18) svc="cache" ;;
        19) svc="registry" ;;
        20) svc="all" ;;
        0) return ;;
        *)
            print_error "无效选择"
            return
            ;;
    esac
    
    case $action in
        start)
            if is_compose_service "$svc"; then
                start_compose_service "$svc"
            else
                for s in $(parse_services "$svc"); do
                    start_service "$s"
                done
            fi
            ;;
        stop)
            if is_compose_service "$svc"; then
                stop_compose_service "$svc"
            else
                for s in $(parse_services "$svc"); do
                    stop_service "$s"
                done
            fi
            ;;
        restart)
            if is_compose_service "$svc"; then
                stop_compose_service "$svc"
                start_compose_service "$svc"
            else
                for s in $(parse_services "$svc"); do
                    stop_service "$s"
                    start_service "$s"
                done
            fi
            ;;
        remove)
            for s in $(parse_services "$svc"); do
                remove_service "$s"
            done
            ;;
        logs)
            if is_compose_service "$svc"; then
                print_error "Compose服务请使用 docker compose logs 查看"
            else
                read -p "显示行数 (默认100): " lines
                lines=${lines:-100}
                logs_service "$svc" "$lines"
            fi
            ;;
        pull)
            pull_images "$(parse_services "$svc")"
            ;;
    esac
}

# 交互式菜单
interactive_menu() {
    while true; do
        echo ""
        echo -e "${CYAN}========== GoFrame Docker Services Manager ==========${NC}"
        echo ""
        echo "  1) 启动服务"
        echo "  2) 停止服务"
        echo "  3) 重启服务"
        echo "  4) 删除服务"
        echo "  5) 查看日志"
        echo "  6) 查看状态"
        echo "  7) 服务信息"
        echo "  8) 清理所有容器"
        echo "  9) 拉取镜像"
        echo "  0) 退出"
        echo ""
        read -p "请选择操作 [0-9]: " choice
        
        case $choice in
            1)
                select_service_menu "start" "启动"
                ;;
            2)
                select_service_menu "stop" "停止"
                ;;
            3)
                select_service_menu "restart" "重启"
                ;;
            4)
                select_service_menu "remove" "删除"
                ;;
            5)
                select_service_menu "logs" "查看日志"
                ;;
            6)
                show_status
                ;;
            7)
                show_service_info
                ;;
            8)
                read -p "确认删除所有goframe容器? [y/N]: " confirm
                if [[ "$confirm" =~ ^[Yy]$ ]]; then
                    clean_all
                fi
                ;;
            9)
                select_service_menu "pull" "拉取镜像"
                ;;
            0)
                echo "再见!"
                exit 0
                ;;
            *)
                print_error "无效选择"
                ;;
        esac
    done
}

# 主函数
main() {
    check_docker
    
    if [ $# -eq 0 ]; then
        interactive_menu
        exit 0
    fi
    
    local command=$1
    local target=$2
    local extra=$3
    
    case $command in
        start)
            if [ -z "$target" ]; then
                print_error "请指定服务名或服务组"
                exit 1
            fi
            if is_compose_service "$target"; then
                start_compose_service "$target"
            else
                for service in $(parse_services "$target"); do
                    start_service "$service"
                done
            fi
            ;;
        stop)
            if [ -z "$target" ]; then
                print_error "请指定服务名或服务组"
                exit 1
            fi
            if is_compose_service "$target"; then
                stop_compose_service "$target"
            else
                for service in $(parse_services "$target"); do
                    stop_service "$service"
                done
            fi
            ;;
        restart)
            if [ -z "$target" ]; then
                print_error "请指定服务名或服务组"
                exit 1
            fi
            if is_compose_service "$target"; then
                stop_compose_service "$target"
                start_compose_service "$target"
            else
                for service in $(parse_services "$target"); do
                    stop_service "$service"
                    start_service "$service"
                done
            fi
            ;;
        remove|rm)
            if [ -z "$target" ]; then
                print_error "请指定服务名或服务组"
                exit 1
            fi
            for service in $(parse_services "$target"); do
                remove_service "$service"
            done
            ;;
        logs)
            if [ -z "$target" ]; then
                print_error "请指定服务名"
                exit 1
            fi
            logs_service "$target" "${extra:-100}"
            ;;
        status|ps)
            show_status
            ;;
        info|list)
            show_service_info
            ;;
        clean)
            clean_all
            ;;
        pull)
            pull_images "$target"
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "未知命令: $command"
            show_help
            exit 1
            ;;
    esac
}

main "$@"

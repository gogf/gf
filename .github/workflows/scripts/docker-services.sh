#!/usr/bin/env bash
#
# GoFrame Docker Services Manager
# For managing Docker services used in local development and testing
#

set -e

# Container name prefix
PREFIX="goframe"

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Service definitions
declare -A SERVICES
declare -A SERVICE_PORTS
declare -A SERVICE_ENVS
declare -A SERVICE_OPTS

# Basic services
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

# Service groups
GROUP_DB="mysql mariadb postgres mssql oracle dm gaussdb clickhouse"
GROUP_CACHE="redis etcd"
GROUP_REGISTRY="polaris zookeeper"
GROUP_ALL="etcd redis mysql mariadb postgres mssql clickhouse polaris oracle dm gaussdb zookeeper"

# Working directories
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
WORKFLOW_DIR="$PROJECT_ROOT/.github/workflows"

# Print colored messages
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

# Check if Docker is available
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed or not in PATH"
        exit 1
    fi
    if ! docker info &> /dev/null; then
        print_error "Docker service is not running"
        exit 1
    fi
}

# Get container name
get_container_name() {
    echo "${PREFIX}-$1"
}

# Start a single service
start_service() {
    local service=$1
    local container_name=$(get_container_name "$service")
    local image="${SERVICES[$service]}"
    local ports="${SERVICE_PORTS[$service]}"
    local envs="${SERVICE_ENVS[$service]}"
    local opts="${SERVICE_OPTS[$service]}"

    if [ -z "$image" ]; then
        print_error "Unknown service: $service"
        return 1
    fi

    # Check if container already exists
    if docker ps -a --format '{{.Names}}' | grep -q "^${container_name}$"; then
        if docker ps --format '{{.Names}}' | grep -q "^${container_name}$"; then
            print_warning "$service is already running"
            return 0
        else
            print_info "Starting existing container $service..."
            docker start "$container_name" > /dev/null
            print_success "$service started"
            return 0
        fi
    fi

    print_info "Starting $service..."
    
    # Build docker run command
    local cmd="docker run -d --name $container_name"
    
    # Add port mappings
    for port in $ports; do
        cmd="$cmd -p $port"
    done
    
    # Add environment variables
    if [ -n "$envs" ]; then
        cmd="$cmd $envs"
    fi
    
    # Add other options
    if [ -n "$opts" ]; then
        cmd="$cmd $opts"
    fi
    
    cmd="$cmd $image"
    
    if eval "$cmd" > /dev/null 2>&1; then
        print_success "$service started (container: $container_name)"
    else
        print_error "Failed to start $service"
        return 1
    fi
}

# Stop a single service
stop_service() {
    local service=$1
    local container_name=$(get_container_name "$service")

    if docker ps --format '{{.Names}}' | grep -q "^${container_name}$"; then
        print_info "Stopping $service..."
        docker stop "$container_name" > /dev/null
        print_success "$service stopped"
    else
        print_warning "$service is not running"
    fi
}

# Remove a single service
remove_service() {
    local service=$1
    local container_name=$(get_container_name "$service")

    if docker ps -a --format '{{.Names}}' | grep -q "^${container_name}$"; then
        print_info "Removing $service..."
        docker rm -f "$container_name" > /dev/null
        print_success "$service removed"
    else
        print_warning "$service container does not exist"
    fi
}

# View service logs
logs_service() {
    local service=$1
    local container_name=$(get_container_name "$service")
    local lines=${2:-100}

    if docker ps -a --format '{{.Names}}' | grep -q "^${container_name}$"; then
        docker logs --tail "$lines" -f "$container_name"
    else
        print_error "$service container does not exist"
        return 1
    fi
}

# Start docker-compose service
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
            print_error "Unknown compose service: $service"
            return 1
            ;;
    esac
    
    if [ -f "$compose_file" ]; then
        print_info "Starting $service (docker-compose)..."
        docker compose -f "$compose_file" up -d
        print_success "$service started"
    else
        print_error "Compose file does not exist: $compose_file"
        return 1
    fi
}

# Stop docker-compose service
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
            print_error "Unknown compose service: $service"
            return 1
            ;;
    esac
    
    if [ -f "$compose_file" ]; then
        print_info "Stopping $service (docker-compose)..."
        docker compose -f "$compose_file" down
        print_success "$service stopped"
    else
        print_error "Compose file does not exist: $compose_file"
        return 1
    fi
}

# Show service status
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

# Show service information
show_service_info() {
    echo ""
    echo -e "${CYAN}========== Available Services ==========${NC}"
    echo ""
    echo -e "${YELLOW}Basic Services (standalone containers):${NC}"
    echo ""
    printf "%-15s %-50s %s\n" "SERVICE" "IMAGE" "PORTS"
    echo "--------------------------------------------------------------------------------"
    
    for service in $GROUP_ALL; do
        printf "%-15s %-50s %s\n" "$service" "${SERVICES[$service]}" "${SERVICE_PORTS[$service]}"
    done
    
    echo ""
    echo -e "${YELLOW}Compose Services (multi-container):${NC}"
    echo "  apollo        - Apollo Config Center (8080, 8070, 8060, 13306)"
    echo "  nacos         - Nacos Registry (8848, 9848, 9555)"
    echo "  redis-cluster - Redis Master-Slave + Sentinel Cluster (6380-6382, 26379-26381)"
    echo "  consul        - Consul Service Discovery (8500, 8600)"
    echo ""
    echo -e "${YELLOW}Service Groups:${NC}"
    echo "  db       - Databases: $GROUP_DB"
    echo "  cache    - Cache: $GROUP_CACHE"
    echo "  registry - Registry: $GROUP_REGISTRY"
    echo "  all      - All basic services"
    echo ""
}

# Show help
show_help() {
    echo ""
    echo -e "${CYAN}GoFrame Docker Services Manager${NC}"
    echo ""
    echo "Usage: $0 <command> [service|group] [options]"
    echo ""
    echo "Commands:"
    echo "  start <service|group>    Start service or service group"
    echo "  stop <service|group>     Stop service or service group"
    echo "  restart <service|group>  Restart service or service group"
    echo "  remove <service|group>   Remove service container"
    echo "  logs <service> [lines]   View service logs (default 100 lines)"
    echo "  status                   Show all service status"
    echo "  info                     Show available service information"
    echo "  clean                    Remove all goframe containers"
    echo "  pull [service]           Pull images"
    echo ""
    echo "Services:"
    echo "  Basic: etcd, redis, mysql, mariadb, postgres, mssql,"
    echo "         clickhouse, polaris, oracle, dm, gaussdb, zookeeper"
    echo "  Compose: apollo, nacos, redis-cluster, consul"
    echo ""
    echo "Service Groups:"
    echo "  db       - All database services"
    echo "  cache    - Cache services (redis, etcd)"
    echo "  registry - Registry services (polaris, zookeeper)"
    echo "  all      - All basic services"
    echo ""
    echo "Examples:"
    echo "  $0 start mysql           # Start MySQL"
    echo "  $0 start db              # Start all databases"
    echo "  $0 start all             # Start all basic services"
    echo "  $0 start apollo          # Start Apollo (compose)"
    echo "  $0 stop all              # Stop all basic services"
    echo "  $0 logs mysql 50         # View MySQL last 50 lines of logs"
    echo "  $0 status                # View service status"
    echo ""
}

# Parse service groups
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

# Check if it's a compose service
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

# Pull images
pull_images() {
    local services=$1
    
    if [ -z "$services" ]; then
        services="$GROUP_ALL"
    fi
    
    for service in $services; do
        if [ -n "${SERVICES[$service]}" ]; then
            print_info "Pulling image: ${SERVICES[$service]}"
            docker pull "${SERVICES[$service]}"
        fi
    done
}

# Clean all goframe containers
clean_all() {
    print_info "Removing all $PREFIX containers..."
    local containers=$(docker ps -a --filter "name=$PREFIX" --format '{{.Names}}')
    
    if [ -n "$containers" ]; then
        for container in $containers; do
            docker rm -f "$container" > /dev/null
            print_success "Removed: $container"
        done
    else
        print_info "No $PREFIX containers found"
    fi
}

# Get service status mark
get_service_status_mark() {
    local service=$1
    local container_name=$(get_container_name "$service")
    
    if docker ps --format '{{.Names}}' 2>/dev/null | grep -q "^${container_name}$"; then
        echo -e "${GREEN}*${NC}"
    else
        echo " "
    fi
}

# Get compose service status mark
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

# Service selection menu
select_service_menu() {
    local action=$1
    local action_name=$2
    
    echo ""
    echo -e "${CYAN}========== Select Service to ${action_name} ==========${NC}"
    
    # Show running status for stop/restart/logs operations
    if [[ "$action" == "stop" || "$action" == "restart" || "$action" == "logs" ]]; then
        echo -e "  (${GREEN}*${NC} indicates running)"
    fi
    echo ""
    echo -e "${YELLOW}Basic Services:${NC}"
    printf "  %b1) etcd         %b2) redis        %b3) mysql\n" \
        "$(get_service_status_mark etcd)" "$(get_service_status_mark redis)" "$(get_service_status_mark mysql)"
    printf "  %b4) mariadb      %b5) postgres     %b6) mssql\n" \
        "$(get_service_status_mark mariadb)" "$(get_service_status_mark postgres)" "$(get_service_status_mark mssql)"
    printf "  %b7) clickhouse   %b8) polaris      %b9) oracle\n" \
        "$(get_service_status_mark clickhouse)" "$(get_service_status_mark polaris)" "$(get_service_status_mark oracle)"
    printf " %b10) dm          %b11) gaussdb     %b12) zookeeper\n" \
        "$(get_service_status_mark dm)" "$(get_service_status_mark gaussdb)" "$(get_service_status_mark zookeeper)"
    echo ""
    echo -e "${YELLOW}Compose Services:${NC}"
    printf " %b13) apollo      %b14) nacos       %b15) redis-cluster\n" \
        "$(get_compose_status_mark apollo)" "$(get_compose_status_mark nacos)" "$(get_compose_status_mark redis-cluster)"
    printf " %b16) consul\n" "$(get_compose_status_mark consul)"
    echo ""
    echo -e "${YELLOW}Service Groups:${NC}"
    echo "  17) db (all databases)    18) cache (cache services)"
    echo "  19) registry (registries) 20) all (all basic services)"
    echo ""
    echo "   0) Back to main menu"
    echo ""
    read -p "Select [0-20]: " svc_choice
    
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
            print_error "Invalid selection"
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
                print_error "For Compose services, please use 'docker compose logs'"
            else
                read -p "Number of lines (default 100): " lines
                lines=${lines:-100}
                logs_service "$svc" "$lines"
            fi
            ;;
        pull)
            pull_images "$(parse_services "$svc")"
            ;;
    esac
}

# Interactive menu
interactive_menu() {
    while true; do
        echo ""
        echo -e "${CYAN}========== GoFrame Docker Services Manager ==========${NC}"
        echo ""
        echo "  1) Start Service"
        echo "  2) Stop Service"
        echo "  3) Restart Service"
        echo "  4) Remove Service"
        echo "  5) View Logs"
        echo "  6) View Status"
        echo "  7) Service Info"
        echo "  8) Clean All Containers"
        echo "  9) Pull Images"
        echo "  0) Exit"
        echo ""
        read -p "Select operation [0-9]: " choice
        
        case $choice in
            1)
                select_service_menu "start" "Start"
                ;;
            2)
                select_service_menu "stop" "Stop"
                ;;
            3)
                select_service_menu "restart" "Restart"
                ;;
            4)
                select_service_menu "remove" "Remove"
                ;;
            5)
                select_service_menu "logs" "View Logs"
                ;;
            6)
                show_status
                ;;
            7)
                show_service_info
                ;;
            8)
                read -p "Confirm removing all goframe containers? [y/N]: " confirm
                if [[ "$confirm" =~ ^[Yy]$ ]]; then
                    clean_all
                fi
                ;;
            9)
                select_service_menu "pull" "Pull Images"
                ;;
            0)
                echo "Goodbye!"
                exit 0
                ;;
            *)
                print_error "Invalid selection"
                ;;
        esac
    done
}

# Main function
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
                print_error "Please specify service name or service group"
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
                print_error "Please specify service name or service group"
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
                print_error "Please specify service name or service group"
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
                print_error "Please specify service name or service group"
                exit 1
            fi
            for service in $(parse_services "$target"); do
                remove_service "$service"
            done
            ;;
        logs)
            if [ -z "$target" ]; then
                print_error "Please specify service name"
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
            print_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

main "$@"

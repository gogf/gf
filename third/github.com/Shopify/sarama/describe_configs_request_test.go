package sarama

import "testing"

var (
	emptyDescribeConfigsRequest = []byte{
		0, 0, 0, 0, // 0 configs
	}

	singleDescribeConfigsRequest = []byte{
		0, 0, 0, 1, // 1 config
		2,                   // a topic
		0, 3, 'f', 'o', 'o', // topic name: foo
		0, 0, 0, 1, //1 config name
		0, 10, // 10 chars
		's', 'e', 'g', 'm', 'e', 'n', 't', '.', 'm', 's',
	}

	doubleDescribeConfigsRequest = []byte{
		0, 0, 0, 2, // 2 configs
		2,                   // a topic
		0, 3, 'f', 'o', 'o', // topic name: foo
		0, 0, 0, 2, //2 config name
		0, 10, // 10 chars
		's', 'e', 'g', 'm', 'e', 'n', 't', '.', 'm', 's',
		0, 12, // 12 chars
		'r', 'e', 't', 'e', 'n', 't', 'i', 'o', 'n', '.', 'm', 's',
		2,                   // a topic
		0, 3, 'b', 'a', 'r', // topic name: foo
		0, 0, 0, 1, // 1 config
		0, 10, // 10 chars
		's', 'e', 'g', 'm', 'e', 'n', 't', '.', 'm', 's',
	}

	singleDescribeConfigsRequestAllConfigs = []byte{
		0, 0, 0, 1, // 1 config
		2,                   // a topic
		0, 3, 'f', 'o', 'o', // topic name: foo
		255, 255, 255, 255, // all configs
	}

	singleDescribeConfigsRequestAllConfigsv1 = []byte{
		0, 0, 0, 1, // 1 config
		2,                   // a topic
		0, 3, 'f', 'o', 'o', // topic name: foo
		255, 255, 255, 255, // no configs
		1, //synoms
	}
)

func TestDescribeConfigsRequestv0(t *testing.T) {
	var request *DescribeConfigsRequest

	request = &DescribeConfigsRequest{
		Version:   0,
		Resources: []*ConfigResource{},
	}
	testRequest(t, "no requests", request, emptyDescribeConfigsRequest)

	configs := []string{"segment.ms"}
	request = &DescribeConfigsRequest{
		Version: 0,
		Resources: []*ConfigResource{
			&ConfigResource{
				Type:        TopicResource,
				Name:        "foo",
				ConfigNames: configs,
			},
		},
	}

	testRequest(t, "one config", request, singleDescribeConfigsRequest)

	request = &DescribeConfigsRequest{
		Version: 0,
		Resources: []*ConfigResource{
			&ConfigResource{
				Type:        TopicResource,
				Name:        "foo",
				ConfigNames: []string{"segment.ms", "retention.ms"},
			},
			&ConfigResource{
				Type:        TopicResource,
				Name:        "bar",
				ConfigNames: []string{"segment.ms"},
			},
		},
	}
	testRequest(t, "two configs", request, doubleDescribeConfigsRequest)

	request = &DescribeConfigsRequest{
		Version: 0,
		Resources: []*ConfigResource{
			&ConfigResource{
				Type: TopicResource,
				Name: "foo",
			},
		},
	}

	testRequest(t, "one topic, all configs", request, singleDescribeConfigsRequestAllConfigs)
}

func TestDescribeConfigsRequestv1(t *testing.T) {
	var request *DescribeConfigsRequest

	request = &DescribeConfigsRequest{
		Version: 1,
		Resources: []*ConfigResource{
			{
				Type: TopicResource,
				Name: "foo",
			},
		},
		IncludeSynonyms: true,
	}

	testRequest(t, "one topic, all configs", request, singleDescribeConfigsRequestAllConfigsv1)
}

package unicometaschedulertest

import (
        "context"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "os"

        framework "github.com/karmada-io/karmada/pkg/scheduler/framework"
)

const Name = "UNICOMetaSchedulerTest"

// UNICOMetaSchedulerTest implements a scoring plugin that assigns a score to a cluster
// based on values read from a JSON file. The JSON is expected to be in the format:
// {
//   "clusterName": {
//       "scoreType1": scoreValue1,
//       "scoreType2": scoreValue2,
//       ...
//   },
//   ...
// }
type UNICOMetaSchedulerTest struct {
        scores     map[string]map[string]int64
        scoringKey string
}

// Ensure UNICOMetaSchedulerTest implements the ScorePlugin and ScoreExtensions interfaces.
var _ framework.ScorePlugin = &UNICOMetaSchedulerTest{}
var _ framework.ScoreExtensions = &UNICOMetaSchedulerTest{}

// Name returns the plugin name.
func (p *UNICOMetaSchedulerTest) Name() string {
        return Name
}

// Score returns the score for the given cluster based on its name and the selected scoring key.
// If the cluster is not present in the JSON mapping or the scoring key is not found,
// it returns a default score of 0.
func (p *UNICOMetaSchedulerTest) Score(ctx context.Context, state *framework.CycleState, cluster *framework.ClusterInfo, nodeName string) (int64, *framework.Status) {
        clusterScores, ok := p.scores[cluster.Name]
        if !ok {
                return 0, framework.NewStatus(framework.Success, "")
        }
        score, ok := clusterScores[p.scoringKey]
        if !ok {
                score = 0
        }
        return score, framework.NewStatus(framework.Success, "")
}

// NormalizeScore allows for normalizing scores across clusters.
// In this example no normalization is applied.
func (p *UNICOMetaSchedulerTest) NormalizeScore(ctx context.Context, state *framework.CycleState, clusters []framework.ClusterInfo, scores framework.NodeScoreList) *framework.Status {
        // No normalization applied in this example.
        return framework.NewStatus(framework.Success, "")
}

// New is the factory function for the plugin.
// It expects the plugin arguments to contain:
// - "jsonPath": the path to the JSON file.
// - "scoringKey": the key to select the score (e.g., "green-power").
func New(args framework.PluginArgs) (framework.Plugin, error) {
        jsonPath, ok := args["jsonPath"].(string)
        if !ok || jsonPath == "" {
                return nil, fmt.Errorf("jsonPath argument is required")
        }
        scoringKey, ok := args["scoringKey"].(string)
        if !ok || scoringKey == "" {
                return nil, fmt.Errorf("scoringKey argument is required")
        }

        // Open and read the JSON file.
        file, err := os.Open(jsonPath)
        if err != nil {
                return nil, fmt.Errorf("failed to open JSON file: %v", err)
        }
        defer file.Close()

        data, err := ioutil.ReadAll(file)
        if err != nil {
                return nil, fmt.Errorf("failed to read JSON file: %v", err)
        }

        var scores map[string]map[string]int64
        if err := json.Unmarshal(data, &scores); err != nil {
                return nil, fmt.Errorf("failed to unmarshal JSON file: %v", err)
        }

        return &UNICOMetaSchedulerTest{
                scores:     scores,
                scoringKey: scoringKey,
        }, nil
}
package position

import (
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"mlib.com/mlog"
)

func GetPositionMap(folder_path string, file_name string /* = "_position.yaml"*/) map[string]int {
	/*
	   Get the mapping from name to index from a YAML file
	   :param folder_path:
	   :param file_name: the YAML file name, default to '_position.yaml'
	   :return: a dict with name as key and index as value
	*/
	if file_name == "" {
		file_name = "_position.yaml"
	}
	position_file_path := filepath.Join(folder_path, file_name)
	filecontent, err := os.ReadFile(position_file_path)
	if err != nil {
		mlog.Errorf("read file(%s) failed:%v", position_file_path, err)
		return nil
	}
	yaml_content := []any{}
	err = yaml.Unmarshal(filecontent, &yaml_content)
	if err != nil {
		mlog.Errorf("yaml.Unmarshal(%s) failed:%v", string(filecontent), err)
		return nil
	}
	positions := []string{}
	for _, item := range yaml_content {
		if item != nil {
			if _, ok := item.(string); ok && strings.TrimSpace(item.(string)) != "" {
				positions = append(positions, strings.TrimSpace(item.(string)))
			}
		}
	}
	positionmap := map[string]int{}
	for idx, name := range positions {
		positionmap[name] = idx
	}
	return positionmap
}
func GetToolPositionMap(folder_path string, file_name string /* = "_position.yaml"*/) map[string]int {
	/*
	   Get the mapping for tools from name to index from a YAML file.
	   :param folder_path:
	   :param file_name: the YAML file name, default to '_position.yaml'
	   :return: a dict with name as key and index as value
	*/
	position_map := GetPositionMap(folder_path, file_name)

	return PinPositionMap(position_map, viper.GetStringSlice("position.tool_pin_list"))

}
func GetProviderPositionMap(folder_path string, file_name string /* = "_position.yaml"*/) map[string]int {
	/*
	   Get the mapping for providers from name to index from a YAML file.
	   :param folder_path:
	   :param file_name: the YAML file name, default to '_position.yaml'
	   :return: a dict with name as key and index as value
	*/
	position_map := GetPositionMap(folder_path, file_name)
	return PinPositionMap(position_map, viper.GetStringSlice("position.provider_pin_list"))
}
func PinPositionMap(original_position_map map[string]int, pin_list []string) map[string]int {
	/*
	   Pin the items in the pin list to the beginning of the position map.
	   Overall logic: exclude > include > pin
	   :param position_map: the position map to be sorted and filtered
	   :param pin_list: the list of pins to be put at the beginning
	   :return: the sorted position map
	*/
	positions := slices.Sorted(maps.Keys(original_position_map))

	// Add pins to position map
	position_map := map[string]int{}
	for idx, name := range pin_list {
		position_map[name] = idx
	}

	// Add remaining positions to position map
	start_idx := len(position_map)
	for _, name := range positions {
		if _, ok := position_map[name]; !ok {
			position_map[name] = start_idx
			start_idx += 1
		}
	}
	return position_map

}

// func is_filtered(
//     include_set: set[str],
//     exclude_set: set[str],
//     data: Any,
//     name_func: Callable[[Any], str],
// ) -> bool{
//     /*
//     Check if the object should be filtered out.
//     Overall logic: exclude > include > pin
//     :param include_set: the set of names to be included
//     :param exclude_set: the set of names to be excluded
//     :param name_func: the function to get the name of the object
//     :param data: the data to be filtered
//     :return: True if the object should be filtered out, False otherwise
//     */
//     if not data:
//         return False
//     if not include_set and not exclude_set:
//         return False

//     name = name_func(data)

//     if name in exclude_set:  // exclude_set is prioritized
//         return True
//     if include_set and name not in include_set:  // filter out only if include_set is not empty
//         return True
//     return False

// }
// func sort_by_position_map(
//     position_map map[string]int,
//     data []any,
//     name_func func[[Any], str],
// ) []any{
//     /*
//     Sort the objects by the position map.
//     If the name of the object is not in the position map, it will be put at the end.
//     :param position_map: the map holding positions in the form of {name: index}
//     :param name_func: the function to get the name of the object
//     :param data: the data to be sorted
//     :return: the sorted objects
//     */
//     if not position_map or not data:
//         return data

//     return sorted(data, key=lambda x: position_map.get(name_func(x), float("inf")))

// }
// func sort_to_dict_by_position_map(
//     position_map map[string]int,
//     data []any,
//     name_func: Callable[[Any], str],
// ) -> OrderedDict[str, Any]{
//     /*
//     Sort the objects into a ordered dict by the position map.
//     If the name of the object is not in the position map, it will be put at the end.
//     :param position_map: the map holding positions in the form of {name: index}
//     :param name_func: the function to get the name of the object
//     :param data: the data to be sorted
//     :return: an OrderedDict with the sorted pairs of name and object
//     */
//     sorted_items = sort_by_position_map(position_map, data, name_func)
//     return OrderedDict([(name_func(item), item) for item in sorted_items])

// }
